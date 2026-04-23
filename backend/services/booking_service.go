package services

import (
	"epbms/config"
	"epbms/models"
	"time"
)

func HasBookingConflict(performerID uint, eventDate, startTime, endTime string) (bool, error) {
	var bookings []models.Booking

	err := config.DB.Where("performer_id = ? AND event_date = ?", performerID, eventDate).Find(&bookings).Error
	if err != nil {
		return false, err
	}

	newStart, err := time.Parse("15:04", startTime)
	if err != nil {
		return false, err
	}

	newEnd, err := time.Parse("15:04", endTime)
	if err != nil {
		return false, err
	}

	for _, booking := range bookings {
		existingStart, err := time.Parse("15:04", booking.StartTime)
		if err != nil {
			return false, err
		}

		existingEnd, err := time.Parse("15:04", booking.EndTime)
		if err != nil {
			return false, err
		}

		if newStart.Before(existingEnd) && newEnd.After(existingStart) {
			return true, nil
		}
	}

	return false, nil
}

func HasBookingConflictExcludingID(bookingID uint, performerID uint, eventDate, startTime, endTime string) (bool, error) {
	var bookings []models.Booking

	err := config.DB.Where("id <> ? AND performer_id = ? AND event_date = ?", bookingID, performerID, eventDate).Find(&bookings).Error
	if err != nil {
		return false, err
	}

	newStart, err := time.Parse("15:04", startTime)
	if err != nil {
		return false, err
	}

	newEnd, err := time.Parse("15:04", endTime)
	if err != nil {
		return false, err
	}

	for _, booking := range bookings {
		existingStart, err := time.Parse("15:04", booking.StartTime)
		if err != nil {
			return false, err
		}

		existingEnd, err := time.Parse("15:04", booking.EndTime)
		if err != nil {
			return false, err
		}

		if newStart.Before(existingEnd) && newEnd.After(existingStart) {
			return true, nil
		}
	}

	return false, nil
}
