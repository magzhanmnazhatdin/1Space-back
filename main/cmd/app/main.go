package main

import (
	"context"
	"log"
	"main/internal/deliveries/handlers"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

// Глобальные переменные
var client *firestore.Client
var firebaseAuth *auth.Client

// Инициализация Firestore и Firebase Auth
func initFirestore() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("../../configs/firebase.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Ошибка подключения к Firebase: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Ошибка создания клиента Firestore: %v", err)
	}

	firebaseAuth, err = app.Auth(ctx)
	if err != nil {
		log.Fatalf("Ошибка инициализации Firebase Auth: %v", err)
	}
}

func main() {
	initFirestore()
	defer client.Close()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handlers.Init(client)

	// Открытые маршруты
	r.GET("/clubs", handlers.GetAllClubs)
	r.GET("/clubs/:id", handlers.GetClubByID)
	r.POST("/auth", handlers.AuthHandler)
	r.GET("/computers", handlers.GetAllComputers)

	// Защищенные маршруты (только проверка аутентификации)
	r.POST("/clubs", handlers.AuthMiddleware(), handlers.CreateClub)
	r.PUT("/clubs/:id", handlers.AuthMiddleware(), handlers.UpdateClub)
	r.DELETE("/clubs/:id", handlers.AuthMiddleware(), handlers.DeleteClub)

	// Маршруты для бронирований
	r.GET("/clubs/:id/computers", handlers.GetClubComputers)
	r.GET("/bookings", handlers.AuthMiddleware(), handlers.GetUserBookings)
	r.POST("/bookings", handlers.AuthMiddleware(), handlers.CreateBooking)
	r.PUT("/bookings/:id/cancel", handlers.AuthMiddleware(), handlers.CancelBooking)
	authRoutes := r.Group("/")
	authRoutes.Use(handlers.AuthMiddleware())
	{
		authRoutes.POST("/clubs/:id/computers", handlers.CreateComputerList)
	}

	r.Run(":8080")
}
