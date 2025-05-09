package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"main/internal/application/usecase"
	"main/internal/domain/entities"
)

// BookingHandler handles HTTP requests for bookings.
type BookingHandler struct {
	bookingUC usecase.BookingUseCase
	clubUC    usecase.ClubUseCase
}

// NewBookingHandler creates a new BookingHandler with injected use cases.
func NewBookingHandler(
	bookingUC usecase.BookingUseCase,
	clubUC usecase.ClubUseCase,
) *BookingHandler {
	return &BookingHandler{
		bookingUC: bookingUC,
		clubUC:    clubUC,
	}
}

func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	uidIf, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, ok := uidIf.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	list, err := h.bookingUC.GetByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = make([]*entities.Booking, 0)
	}
	c.JSON(http.StatusOK, list)
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	// pull UID out of the context (set by AuthMiddleware)
	uidIf, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, _ := uidIf.(string)

	// bind request
	var req struct {
		ClubID    string    `json:"club_id"`
		PCNumber  int       `json:"pc_number"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// fetch club to get price per hour
	club, err := h.clubUC.GetByID(c.Request.Context(), req.ClubID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid club ID"})
		return
	}

	// calculate duration and total price
	durationHours := req.EndTime.Sub(req.StartTime).Hours()
	totalPrice := durationHours * club.PricePerHour

	// assemble booking entity
	booking := &entities.Booking{
		ClubID:     req.ClubID,
		UserID:     userID,
		PCNumber:   req.PCNumber,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		TotalPrice: totalPrice,
	}

	// create booking
	if err := h.bookingUC.Create(c.Request.Context(), booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	id := c.Param("id")
	if err := h.bookingUC.Cancel(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
