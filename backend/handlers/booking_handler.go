package handlers

import (
	"epbms/config"
	"epbms/models"
	"epbms/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBookings(c *gin.Context) {
	roleValue, _ := c.Get("role")
	userIDValue, _ := c.Get("user_id")

	role := roleValue.(string)
	userID := userIDValue.(uint)

	var bookings []models.Booking
	var err error

	if role == "admin" {
		err = config.DB.Preload("Performer").Find(&bookings).Error
	} else if role == "organizer" {
		err = config.DB.Preload("Performer").Where("organizer_id = ?", userID).Find(&bookings).Error
	} else if role == "performer" {
		var performer models.Performer
		err = config.DB.Where("user_id = ?", userID).First(&performer).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "performer profile not found"})
			return
		}

		err = config.DB.Preload("Performer").Where("performer_id = ?", performer.ID).Find(&bookings).Error
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get bookings",
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func GetBookingByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	roleValue, _ := c.Get("role")
	userIDValue, _ := c.Get("user_id")

	role := roleValue.(string)
	userID := userIDValue.(uint)

	var booking models.Booking
	err = config.DB.Preload("Performer").First(&booking, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	if role == "organizer" && booking.OrganizerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if role == "performer" {
		var performer models.Performer
		err = config.DB.Where("user_id = ?", userID).First(&performer).Error
		if err != nil || booking.PerformerID != performer.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	}

	c.JSON(http.StatusOK, booking)
}

func GetBookingsByDate(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date query parameter is required",
		})
		return
	}

	roleValue, _ := c.Get("role")
	userIDValue, _ := c.Get("user_id")

	role := roleValue.(string)
	userID := userIDValue.(uint)

	var bookings []models.Booking
	var err error

	if role == "admin" {
		err = config.DB.Preload("Performer").Where("event_date = ?", date).Find(&bookings).Error
	} else if role == "organizer" {
		err = config.DB.Preload("Performer").Where("event_date = ? AND organizer_id = ?", date, userID).Find(&bookings).Error
	} else if role == "performer" {
		var performer models.Performer
		err = config.DB.Where("user_id = ?", userID).First(&performer).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "performer profile not found"})
			return
		}

		err = config.DB.Preload("Performer").Where("event_date = ? AND performer_id = ?", date, performer.ID).Find(&bookings).Error
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get bookings by date",
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func GetPerformerSchedule(c *gin.Context) {
	performerIDParam := c.Param("id")
	date := c.Query("date")

	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date query parameter is required",
		})
		return
	}

	performerID, err := strconv.Atoi(performerIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid performer id",
		})
		return
	}

	var performer models.Performer
	err = config.DB.First(&performer, performerID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "performer not found",
		})
		return
	}

	var bookings []models.Booking
	err = config.DB.Preload("Performer").
		Where("performer_id = ? AND event_date = ?", performerID, date).
		Order("start_time asc").
		Find(&bookings).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get performer schedule",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"performer": performer,
		"date":      date,
		"bookings":  bookings,
	})
}

func CreateBooking(c *gin.Context) {
	roleValue, _ := c.Get("role")
	userIDValue, _ := c.Get("user_id")

	role := roleValue.(string)
	userID := userIDValue.(uint)

	if role != "admin" && role != "organizer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin or organizer can create bookings"})
		return
	}

	var booking models.Booking
	err := c.ShouldBindJSON(&booking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if role == "organizer" {
		booking.OrganizerID = userID
	}

	var performer models.Performer
	err = config.DB.First(&performer, booking.PerformerID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "performer not found",
		})
		return
	}

	conflict, err := services.HasBookingConflict(
		booking.PerformerID,
		booking.EventDate,
		booking.StartTime,
		booking.EndTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check booking conflict",
		})
		return
	}

	if conflict {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "booking time conflicts with an existing booking",
		})
		return
	}

	err = config.DB.Create(&booking).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = config.DB.Preload("Performer").First(&booking, booking.ID).Error
	if err != nil {
		c.JSON(http.StatusCreated, booking)
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func UpdateBooking(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	roleValue, _ := c.Get("role")
	userIDValue, _ := c.Get("user_id")

	role := roleValue.(string)
	userID := userIDValue.(uint)

	var booking models.Booking
	err = config.DB.First(&booking, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	if role == "organizer" && booking.OrganizerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if role == "performer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "performers cannot update bookings"})
		return
	}

	var updatedData models.Booking
	err = c.ShouldBindJSON(&updatedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var performer models.Performer
	err = config.DB.First(&performer, updatedData.PerformerID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "performer not found"})
		return
	}

	conflict, err := services.HasBookingConflictExcludingID(
		uint(id),
		updatedData.PerformerID,
		updatedData.EventDate,
		updatedData.StartTime,
		updatedData.EndTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check booking conflict"})
		return
	}

	if conflict {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking time conflicts with an existing booking"})
		return
	}

	booking.PerformerID = updatedData.PerformerID
	booking.ClientName = updatedData.ClientName
	booking.EventDate = updatedData.EventDate
	booking.StartTime = updatedData.StartTime
	booking.EndTime = updatedData.EndTime
	booking.Status = updatedData.Status

	if role == "admin" && updatedData.OrganizerID != 0 {
		booking.OrganizerID = updatedData.OrganizerID
	}

	err = config.DB.Save(&booking).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update booking"})
		return
	}

	err = config.DB.Preload("Performer").First(&booking, booking.ID).Error
	if err != nil {
		c.JSON(http.StatusOK, booking)
		return
	}

	c.JSON(http.StatusOK, booking)
}

func DeleteBooking(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	roleValue, _ := c.Get("role")
	userIDValue, _ := c.Get("user_id")

	role := roleValue.(string)
	userID := userIDValue.(uint)

	var booking models.Booking
	err = config.DB.First(&booking, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	if role == "organizer" && booking.OrganizerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if role == "performer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "performers cannot delete bookings"})
		return
	}

	err = config.DB.Delete(&booking).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "booking deleted successfully",
	})
}
