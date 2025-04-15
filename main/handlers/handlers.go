package handlers

import (
	"log"
	"net/http"

	"main/models"
	"main/services"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"firebase.google.com/go/v4/auth" // Импортируем для работы с Firebase Auth
)

type Handler struct {
	clubService     services.ClubService
	bookingService  services.BookingService
	computerService services.ComputerService
	authClient      *auth.Client // Типизируем как *auth.Client
}

func NewHandler(clubService services.ClubService, bookingService services.BookingService, computerService services.ComputerService, authClient *auth.Client) *Handler {
	return &Handler{
		clubService:     clubService,
		bookingService:  bookingService,
		computerService: computerService,
		authClient:      authClient,
	}
}

func (h *Handler) GetAllClubs(c *gin.Context) {
	clubs, err := h.clubService.GetAllClubs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clubs)
}

func (h *Handler) GetClubByID(c *gin.Context) {
	id := c.Param("id")
	club, err := h.clubService.GetClubByID(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, club)
}

func (h *Handler) CreateClub(c *gin.Context) {
	var club models.ComputerClub
	if err := c.ShouldBindJSON(&club); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.clubService.CreateClub(c.Request.Context(), &club); err != nil {
		if status.Code(err) == codes.InvalidArgument {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Club created", "id": club.ID})
}

func (h *Handler) UpdateClub(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Club ID is required"})
		return
	}

	var club models.ComputerClub
	if err := c.ShouldBindJSON(&club); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Устанавливаем ID из параметра маршрута, чтобы избежать изменения ID
	club.ID = id

	if err := h.clubService.UpdateClub(c.Request.Context(), id, &club); err != nil {
		if status.Code(err) == codes.InvalidArgument {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update club: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Club updated successfully",
		"club":    club,
	})
}

func (h *Handler) DeleteClub(c *gin.Context) {
	id := c.Param("id")
	if err := h.clubService.DeleteClub(c.Request.Context(), id); err != nil {
		if status.Code(err) == codes.InvalidArgument {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Club deleted"})
}

func (h *Handler) GetUserBookings(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	bookings, err := h.bookingService.GetUserBookings(c.Request.Context(), uid.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *Handler) CreateBooking(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input services.BookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.bookingService.CreateBooking(c.Request.Context(), uid.(string), &input)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument || status.Code(err) == codes.FailedPrecondition {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, booking)
}

func (h *Handler) CancelBooking(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	bookingID := c.Param("id")
	if err := h.bookingService.CancelBooking(c.Request.Context(), uid.(string), bookingID); err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		} else if status.Code(err) == codes.InvalidArgument || status.Code(err) == codes.FailedPrecondition {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if status.Code(err) == codes.PermissionDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled"})
}

func (h *Handler) GetAllComputers(c *gin.Context) {
	computers, err := h.computerService.GetAllComputers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, computers)
}

func (h *Handler) GetClubComputers(c *gin.Context) {
	clubID := c.Param("id")
	computers, err := h.computerService.GetComputersByClubID(c.Request.Context(), clubID)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, computers)
}

func (h *Handler) CreateComputerList(c *gin.Context) {
	clubID := c.Param("id")
	var computers []models.Computer
	if err := c.ShouldBindJSON(&computers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdComputers, err := h.computerService.CreateComputers(c.Request.Context(), clubID, computers)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument || status.Code(err) == codes.FailedPrecondition {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":   "Computers created",
		"club_id":   clubID,
		"computers": createdComputers,
	})
}

func (h *Handler) AuthHandler(c *gin.Context) {
	// Извлекаем токен из заголовка Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Printf("AuthHandler: Authorization header is missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Проверяем, что заголовок имеет формат "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		log.Printf("AuthHandler: Invalid Authorization header format: %s", authHeader)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	// Извлекаем сам токен
	token := authHeader[len(bearerPrefix):]
	if token == "" {
		log.Printf("AuthHandler: Token is empty")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	// Проверяем токен с помощью Firebase Admin SDK
	ctx := c.Request.Context()
	decodedToken, err := h.authClient.VerifyIDToken(ctx, token)
	if err != nil {
		log.Printf("AuthHandler: Token verification failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token: " + err.Error()})
		return
	}

	// Токен валиден, извлекаем uid
	uid := decodedToken.UID
	if uid == "" {
		log.Printf("AuthHandler: UID extraction failed: %v", decodedToken)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract UID from token"})
		return
	}

	// Логируем успешную аутентификацию
	log.Printf("AuthHandler: User authenticated: %s", uid)

	// Возвращаем uid в ответе
	c.JSON(http.StatusOK, gin.H{"uid": uid})
}
