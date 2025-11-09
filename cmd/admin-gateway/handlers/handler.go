package handlers

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	"booker/cmd/admin-gateway/config"
	"booker/cmd/admin-gateway/middleware"
	bookingpb "booker/pkg/proto/booking"
	venuepb "booker/pkg/proto/venue"
	"booker/pkg/redis"
)

type Handler struct {
	venueClient   venuepb.VenueServiceClient
	bookingClient bookingpb.BookingServiceClient
	redisClient   *redis.Client
	cfg           *config.Config
}

func New(venueConn, bookingConn *grpc.ClientConn, redisClient *redis.Client, cfg *config.Config) *Handler {
	return &Handler{
		venueClient:   venuepb.NewVenueServiceClient(venueConn),
		bookingClient: bookingpb.NewBookingServiceClient(bookingConn),
		redisClient:   redisClient,
		cfg:           cfg,
	}
}

func (h *Handler) SetupRoutes(mw *middleware.Middleware) *echo.Echo {
	e := echo.New()

	// Add metrics middleware FIRST to log all requests before Echo Logger
	e.Use(middleware.MetricsMiddleware("admin-gateway"))

	// Setup basic middleware (Logger, Recover, CORS)
	mw.SetupMiddleware(e)

	// Metrics endpoint
	e.GET("/metrics", h.Metrics)

	// API routes - register BEFORE static files to avoid conflicts
	api := e.Group("/api/v1")

	// API info endpoint
	e.GET("/api", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"service": "Admin Gateway",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"auth":         "/api/v1/auth/login",
				"venues":       "/api/v1/venues",
				"bookings":     "/api/v1/bookings",
				"availability": "/api/v1/availability/check",
				"websocket":    "/api/v1/ws",
			},
		})
	})

	// Auth
	api.POST("/auth/login", h.Login)
	api.POST("/auth/refresh", h.RefreshToken)

	// Protected routes
	protected := api.Group("", mw.AuthMiddleware)

	// Venues
	protected.GET("/venues", h.ListVenues)
	protected.GET("/venues/:id", h.GetVenue)
	protected.POST("/venues", h.CreateVenue)
	protected.PUT("/venues/:id", h.UpdateVenue)
	protected.DELETE("/venues/:id", h.DeleteVenue)

	// Rooms
	protected.GET("/venues/:venueId/rooms", h.ListRooms)
	protected.GET("/rooms/:id", h.GetRoom)
	protected.POST("/venues/:venueId/rooms", h.CreateRoom)
	protected.PUT("/rooms/:id", h.UpdateRoom)
	protected.DELETE("/rooms/:id", h.DeleteRoom)

	// Tables
	protected.GET("/rooms/:roomId/tables", h.ListTables)
	protected.GET("/tables/:id", h.GetTable)
	protected.POST("/rooms/:roomId/tables", h.CreateTable)
	protected.PUT("/tables/:id", h.UpdateTable)
	protected.DELETE("/tables/:id", h.DeleteTable)

	// Schedule
	protected.GET("/venues/:venueId/schedule", h.GetOpeningHours)
	protected.POST("/venues/:venueId/schedule", h.SetOpeningHours)
	protected.POST("/venues/:venueId/special-hours", h.SetSpecialHours)

	// Bookings
	protected.GET("/bookings", h.ListBookings)
	protected.GET("/bookings/:id", h.GetBooking)
	protected.POST("/bookings", h.CreateBooking)
	protected.POST("/bookings/:id/confirm", h.ConfirmBooking)
	protected.POST("/bookings/:id/cancel", h.CancelBooking)
	protected.POST("/bookings/:id/seat", h.MarkSeated)
	protected.POST("/bookings/:id/finish", h.MarkFinished)
	protected.POST("/bookings/:id/no-show", h.MarkNoShow)

	// Availability
	protected.POST("/availability/check", h.CheckAvailability)

	// WebSocket
	protected.GET("/ws", h.WebSocket)

	// Static files - serve frontend (register AFTER API routes to avoid conflicts)
	// This will serve index.html on root path and all other static assets
	e.Static("/", "web/dist")

	return e
}
