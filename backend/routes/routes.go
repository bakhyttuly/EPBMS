package routes

import (
	"epbms/internal/domain"
	"epbms/internal/handler"
	"epbms/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all API routes with their respective handlers and middleware.
func SetupRoutes(
	r *gin.Engine,
	authH *handler.AuthHandler,
	performerH *handler.PerformerHandler,
	bookingH *handler.BookingHandler,
	adminH *handler.AdminHandler,
) {
	api := r.Group("/api/v1")

	// --- Public routes (no auth required) ---
	auth := api.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
	}

	// --- Protected routes (JWT required) ---
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		// Performers — readable by all authenticated users; writable by admin
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

		// Bookings — role-scoped visibility enforced in the service layer
		bookings := protected.Group("/bookings")
		{
			bookings.GET("", bookingH.GetAll)
			bookings.GET("/:id", bookingH.GetByID)
			bookings.POST("",
				middleware.RequireRoles(domain.RoleClient),
				bookingH.Create,
			)
		}

		// Admin-only routes
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRoles(domain.RoleAdmin))
		{
			admin.GET("/stats", adminH.GetStats)
			admin.PUT("/bookings/:id/status", bookingH.UpdateStatus)
			admin.DELETE("/bookings/:id", bookingH.Delete)
		}
	}
}
