package domain

import "time"

// BookingStatus represents the lifecycle state of a booking.
type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusRejected  BookingStatus = "rejected"
	StatusCompleted BookingStatus = "completed"
)

// Booking represents a booking request made by a client for a performer.
type Booking struct {
	ID          uint          `gorm:"primaryKey"                                          json:"id"`
	PerformerID uint          `gorm:"not null;index"                                      json:"performer_id"`
	Performer   Performer     `gorm:"foreignKey:PerformerID;constraint:OnDelete:RESTRICT"  json:"performer,omitempty"`
	ClientID    uint          `gorm:"not null;index"                                      json:"client_id"`
	Client      User          `gorm:"foreignKey:ClientID;constraint:OnDelete:RESTRICT"     json:"-"`
	EventDate   string        `gorm:"type:date;not null;index"                            json:"event_date"`
	StartTime   string        `gorm:"type:time;not null"                                  json:"start_time"`
	EndTime     string        `gorm:"type:time;not null"                                  json:"end_time"`
	Status      BookingStatus `gorm:"type:varchar(20);not null;default:'pending';index"   json:"status"`
	Notes       string        `gorm:"type:text"                                           json:"notes"`
	ApprovedBy  *uint         `gorm:"index"                                               json:"approved_by,omitempty"`
	ApprovedAt  *time.Time    `                                                           json:"approved_at,omitempty"`
	CreatedAt   time.Time     `                                                           json:"created_at"`
	UpdatedAt   time.Time     `                                                           json:"updated_at"`
}

// CreateBookingRequest is the payload for a client creating a booking request.
type CreateBookingRequest struct {
	PerformerID uint   `json:"performer_id" binding:"required"`
	EventDate   string `json:"event_date"   binding:"required"`
	StartTime   string `json:"start_time"   binding:"required"`
	EndTime     string `json:"end_time"     binding:"required"`
	Notes       string `json:"notes"`
}

// UpdateBookingStatusRequest is the payload for an admin approving or rejecting a booking.
type UpdateBookingStatusRequest struct {
	Status BookingStatus `json:"status" binding:"required,oneof=confirmed rejected completed"`
}

// BookingFilter holds query parameters for listing bookings.
type BookingFilter struct {
	PerformerID uint
	ClientID    uint
	Status      BookingStatus
	EventDate   string
	Page        int
	PageSize    int
}
