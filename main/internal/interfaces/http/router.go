package http

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"main/internal/config"
	"main/internal/infrastructure/stripeclient"
	"main/internal/interfaces/http/handler"
	"main/internal/interfaces/http/middleware"
)

func NewRouter(
	clubH *handler.ClubHandler,
	compH *handler.ComputerHandler,
	bookH *handler.BookingHandler,
	authH *handler.AuthHandler,
	paymentH *handler.PaymentHandler,
	authClient *auth.Client,
) *gin.Engine {
	// load config
	config.Init()

	// init stripeclient
	stripeclient.Init(config.Cfg.Stripe.SecretKey)

	r := gin.Default()

	// Public routes
	r.POST("/auth", authH.Auth)
	r.GET("/clubs", clubH.GetAllClubs)
	r.GET("/clubs/:id", clubH.GetClubByID)
	r.GET("/computers", compH.GetAllComputers)
	r.GET("/clubs/:id/computers", compH.GetClubComputers)
	r.POST("/payments/create", paymentH.CreateIntent)
	r.POST("/webhook", paymentH.Webhook)

	// Protected routes
	protected := r.Group("/", middleware.AuthMiddleware(authClient))
	{
		protected.POST("/clubs", clubH.CreateClub)
		protected.PUT("/clubs/:id", clubH.UpdateClub)
		protected.DELETE("/clubs/:id", clubH.DeleteClub)

		protected.GET("/bookings", bookH.GetUserBookings)
		protected.POST("/bookings", bookH.CreateBooking)
		protected.PUT("/bookings/:id/cancel", bookH.CancelBooking)

		protected.POST("/clubs/:id/computers", compH.CreateComputerList)
	}
	return r
}
