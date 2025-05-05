package handler

import (
	"github.com/gin-gonic/gin"
	"main/internal/application/usecase"
	"main/internal/domain/entities"
	"net/http"
	"time"
)

// BookingHandler handles HTTP requests for bookings.
type BookingHandler struct {
	uc usecase.BookingUseCase
}

// NewBookingHandler creates a new BookingHandler with injected use case.
func NewBookingHandler(uc usecase.BookingUseCase) *BookingHandler {
	return &BookingHandler{uc: uc}
}

func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	bookings, err := h.uc.GetByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
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
	booking := &entities.Booking{
		ClubID:     req.ClubID,
		UserID:     userID,
		PCNumber:   req.PCNumber,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		TotalPrice: req.EndTime.Sub(req.StartTime).Hours() * 0, // calculate price
	}
	if err := h.uc.Create(c.Request.Context(), booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	id := c.Param("id")
	if err := h.uc.Cancel(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
