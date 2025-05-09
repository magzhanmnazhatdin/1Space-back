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
	userH *handler.UserHandler,
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
	protected := r.Group("/",
		middleware.AuthMiddleware(authClient),
		middleware.RequireRole("user"))
	{
		protected.GET("/bookings", bookH.GetUserBookings)
		protected.POST("/bookings", bookH.CreateBooking)
		protected.PUT("/bookings/:id/cancel", bookH.CancelBooking)
	}

	// маршруты для менеджеров (manager + admin)
	manager := r.Group("/manager",
		middleware.AuthMiddleware(authClient),
		middleware.RequireRole("manager", "admin"),
	)
	{
		manager.POST("/clubs", clubH.CreateClub)
		manager.PUT("/clubs/:id", clubH.UpdateClub)
		manager.DELETE("/clubs/:id", clubH.DeleteClub)
		manager.POST("/clubs/:id/computers", compH.CreateComputerList)
		manager.PUT("/computers/:id", compH.UpdateComputer)
		manager.DELETE("/computers/:id", compH.DeleteComputer)
	}

	// маршруты для админов (только admin)
	admin := r.Group("/admin",
		middleware.AuthMiddleware(authClient),
		middleware.RequireRole("admin"),
	)
	{
		// управление пользователями, в т.ч. смена роли
		admin.PUT("/users/:id/role", userH.ChangeRole)
		// сюда можно добавить ещё endpoints для админа
	}

	return r
}
