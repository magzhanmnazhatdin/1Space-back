package main

import (
	"context"
	"log"
	"main/internal/config"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"main/internal/application/usecase"
	fsrepo "main/internal/infrastructure/firestore"
	"main/internal/interfaces/http"
	"main/internal/interfaces/http/handler"
)

func main() {
	// Initialize Firebase App
	opt := option.WithCredentialsFile(config.Cfg.Firebase.CredentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase: %v", err)
	}

	// Initialize Firestore
	fsClient, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing firestore: %v", err)
	}
	defer fsClient.Close()

	// Initialize Auth
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error initializing auth: %v", err)
	}

	// Repositories
	clubRepo := fsrepo.NewClubRepoFS(fsClient)
	compRepo := fsrepo.NewComputerRepoFS(fsClient)
	bookRepo := fsrepo.NewBookingRepoFS(fsClient)

	// Use Cases
	clubUC := usecase.NewClubUseCase(clubRepo)
	compUC := usecase.NewComputerUseCase(compRepo)
	bookUC := usecase.NewBookingUseCase(bookRepo, compRepo)
	paymentUC := usecase.NewPaymentUseCase()

	// Handlers
	clubH := handler.NewClubHandler(clubUC)
	compH := handler.NewComputerHandler(compUC, clubUC)
	bookH := handler.NewBookingHandler(bookUC, clubUC)
	authH := handler.NewAuthHandler(authClient)
	paymentH := handler.NewPaymentHandler(paymentUC)
	userH := handler.NewUserHandler(authClient)

	// Router setup
	router := http.NewRouter(clubH, compH, bookH, authH, paymentH, userH, authClient)
	log.Fatal(router.Run(":8080"))
}
