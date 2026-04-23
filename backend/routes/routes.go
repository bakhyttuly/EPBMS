package routes

import (
	"epbms/internal/domain"
	"epbms/internal/handler"
	"epbms/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all API routes and HTML page routes.
func SetupRoutes(
	r *gin.Engine,
	authH *handler.AuthHandler,
	performerH *handler.PerformerHandler,
	bookingH *handler.BookingHandler,
	adminH *handler.AdminHandler,
	pageH *handler.PageHandler,
) {
	// --- Static Files ---
	r.Static("/static", "../frontend/static")

	// --- HTML Page Routes ---
	r.GET("/", pageH.ShowLoginPage)
	r.GET("/login", pageH.ShowLoginPage)
	r.GET("/register", pageH.ShowRegisterPage)
	
	// Pages that require authentication (handled by JS on the client side)
	r.GET("/dashboard-page", pageH.ShowDashboardPage)
	r.GET("/performers-page", pageH.ShowPerformersPage)
	r.GET("/bookings-page", pageH.ShowBookingsPage)
	r.GET("/calendar-page", pageH.ShowCalendarPage)
	r.GET("/my-schedule-page", pageH.ShowMySchedulePage)

	// --- API Routes ---
	api := r.Group("/api/v1")

	// Public API routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
	}

	// Protected API routes
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		// Performers
		performers := protected.Group("/performers")
		{
			performers.GET("", performerH.GetAll)
			performers.GET("/:id", performerH.GetByID)
			performers.POST("",
				middleware.RequireRoles(domain.RoleAdmin),
				performerH.Create,
			)
			performers.PUT("/:id",
				middleware.RequireRoles(domain.RoleAdmin, domain.RolePerformer),
				performerH.Update,
			)
			performers.DELETE("/:id",
				middleware.RequireRoles(domain.RoleAdmin),
				performerH.Delete,
			)
		}

		// Bookings
		bookings := protected.Group("/bookings")
		{
			bookings.GET("", bookingH.GetAll)
			bookings.GET("/:id", bookingH.GetByID)
			bookings.POST("",
				middleware.RequireRoles(domain.RoleClient),
				bookingH.Create,
			)
		}

		// Admin-only
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRoles(domain.RoleAdmin))
		{
			admin.GET("/stats", adminH.GetStats)
			admin.PUT("/bookings/:id/status", bookingH.UpdateStatus)
			admin.DELETE("/bookings/:id", bookingH.Delete)
		}
	}
}
