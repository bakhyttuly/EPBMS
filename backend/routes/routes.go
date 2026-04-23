package routes

import (
	"epbms/handlers"
	"epbms/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", handlers.ShowLoginPage)
	r.GET("/register", handlers.ShowRegisterPage)

	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)
	r.GET("/logout", handlers.Logout)

	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/me", handlers.GetCurrentUser)

		protected.GET("/dashboard-page", middleware.RoleRequired("admin", "organizer"), handlers.ShowDashboardPage)
		protected.GET("/performers-page", middleware.RoleRequired("admin"), handlers.ShowPerformersPage)
		protected.GET("/bookings-page", middleware.RoleRequired("admin", "organizer"), handlers.ShowBookingsPage)
		protected.GET("/calendar-page", middleware.RoleRequired("admin", "organizer"), handlers.ShowCalendarPage)
		protected.GET("/my-schedule-page", middleware.RoleRequired("performer"), handlers.ShowMySchedulePage)

		protected.GET("/dashboard/stats", handlers.GetDashboardStats)

		protected.GET("/performers", handlers.GetPerformers)
		protected.GET("/performers/:id", handlers.GetPerformerByID)
		protected.POST("/performers", handlers.CreatePerformer)
		protected.PUT("/performers/:id", handlers.UpdatePerformer)
		protected.DELETE("/performers/:id", handlers.DeletePerformer)

		protected.GET("/bookings", handlers.GetBookings)
		protected.GET("/bookings/:id", handlers.GetBookingByID)
		protected.GET("/bookings/by-date", handlers.GetBookingsByDate)
		protected.POST("/bookings", handlers.CreateBooking)
		protected.PUT("/bookings/:id", handlers.UpdateBooking)
		protected.DELETE("/bookings/:id", handlers.DeleteBooking)

		protected.GET("/performers/:id/schedule", handlers.GetPerformerSchedule)
	}
}
