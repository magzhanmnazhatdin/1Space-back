package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"main/handlers"
	"main/middleware"
	"main/repositories"
	"main/services"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func initFirestore() (*firestore.Client, *auth.Client, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("firebase.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, nil, err
	}

	firebaseAuth, err := app.Auth(ctx)
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, firebaseAuth, nil
}

func setupHealthCheck(r *gin.Engine, firestoreClient *firestore.Client, authClient *auth.Client) {
	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		response := HealthCheckResponse{
			Status:    "available",
			Version:   "1.0.0",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Services:  make(map[string]string),
		}

		if _, err := firestoreClient.Collection("test").Limit(1).Documents(ctx).Next(); err != nil {
			response.Status = "degraded"
			response.Services["firestore"] = "unavailable: " + err.Error()
		} else {
			response.Services["firestore"] = "connected"
		}

		if _, err := authClient.GetUser(ctx, "test"); err != nil {
			if !auth.IsUserNotFound(err) {
				response.Status = "degraded"
				response.Services["firebase_auth"] = "unavailable: " + err.Error()
			} else {
				response.Services["firebase_auth"] = "connected"
			}
		} else {
			response.Services["firebase_auth"] = "connected"
		}

		statusCode := http.StatusOK
		if response.Status == "degraded" {
			statusCode = http.StatusPartialContent
		}

		c.JSON(statusCode, response)
	})
}

func main() {
	client, firebaseAuth, err := initFirestore()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	defer client.Close()

	// Инициализация репозиториев
	clubRepo := repositories.NewClubRepository(client)
	bookingRepo := repositories.NewBookingRepository(client)
	computerRepo := repositories.NewComputerRepository(client)

	// Инициализация сервисов
	clubService := services.NewClubService(clubRepo)
	bookingService := services.NewBookingService(bookingRepo, clubRepo, computerRepo, client)
	computerService := services.NewComputerService(computerRepo, clubRepo)

	// Инициализация обработчиков
	// Обновляем вызов NewHandler, так как он больше не возвращает ошибку
	handler := handlers.NewHandler(clubService, bookingService, computerService, firebaseAuth)

	// Настройка маршрутов
	r := gin.Default()
	setupHealthCheck(r, client, firebaseAuth)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Открытые маршруты
	r.GET("/clubs", handler.GetAllClubs)
	r.GET("/clubs/:id", handler.GetClubByID)
	r.POST("/auth", handler.AuthHandler)
	r.GET("/computers", handler.GetAllComputers)

	// Защищенные маршруты
	r.POST("/clubs", middleware.AuthMiddleware(firebaseAuth), handler.CreateClub)
	r.PUT("/clubs/:id", middleware.AuthMiddleware(firebaseAuth), handler.UpdateClub)
	r.DELETE("/clubs/:id", middleware.AuthMiddleware(firebaseAuth), handler.DeleteClub)

	// Маршруты для бронирований
	r.GET("/clubs/:id/computers", handler.GetClubComputers)
	r.GET("/bookings", middleware.AuthMiddleware(firebaseAuth), handler.GetUserBookings)
	r.POST("/bookings", middleware.AuthMiddleware(firebaseAuth), handler.CreateBooking)
	r.PUT("/bookings/:id/cancel", middleware.AuthMiddleware(firebaseAuth), handler.CancelBooking)

	// Защищенные маршруты для компьютеров
	authRoutes := r.Group("/").Use(middleware.AuthMiddleware(firebaseAuth))
	{
		authRoutes.POST("/clubs/:id/computers", handler.CreateComputerList)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
