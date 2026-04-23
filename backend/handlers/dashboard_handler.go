package handlers

import (
	"epbms/config"
	"epbms/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardStats struct {
	TotalBookings     int64 `json:"total_bookings"`
	ActiveBookings    int64 `json:"active_bookings"`
	CompletedBookings int64 `json:"completed_bookings"`
	PerformersCount   int64 `json:"performers_count"`
}

func GetDashboardStats(c *gin.Context) {
	var totalBookings int64
	var activeBookings int64
	var completedBookings int64
	var performersCount int64

	if err := config.DB.Model(&models.Booking{}).Count(&totalBookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to count total bookings",
		})
		return
	}

	if err := config.DB.Model(&models.Booking{}).Where("status = ?", "active").Count(&activeBookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to count active bookings",
		})
		return
	}

	if err := config.DB.Model(&models.Booking{}).Where("status = ?", "completed").Count(&completedBookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to count completed bookings",
		})
		return
	}

	if err := config.DB.Model(&models.Performer{}).Count(&performersCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to count performers",
		})
		return
	}

	stats := DashboardStats{
		TotalBookings:     totalBookings,
		ActiveBookings:    activeBookings,
		CompletedBookings: completedBookings,
		PerformersCount:   performersCount,
	}

	c.JSON(http.StatusOK, stats)
}
