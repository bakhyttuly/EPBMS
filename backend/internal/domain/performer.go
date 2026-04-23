package domain

import "time"

// Performer represents an artist or service provider available for booking.
type Performer struct {
	ID          uint      `gorm:"primaryKey"                    json:"id"`
	UserID      uint      `gorm:"uniqueIndex;not null"          json:"user_id"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Name        string    `gorm:"not null"                      json:"name"`
	Category    string    `gorm:"not null"                      json:"category"`
	Price       float64   `gorm:"not null;default:0"            json:"price"`
	Description string    `gorm:"type:text"                     json:"description"`
	CreatedAt   time.Time `                                     json:"created_at"`
	UpdatedAt   time.Time `                                     json:"updated_at"`
}

// CreatePerformerRequest is the payload for creating a performer profile.
type CreatePerformerRequest struct {
	Name        string  `json:"name"        binding:"required,min=2,max=100"`
	Category    string  `json:"category"    binding:"required"`
	Price       float64 `json:"price"       binding:"required,gte=0"`
	Description string  `json:"description"`
}

// UpdatePerformerRequest is the payload for updating a performer profile.
type UpdatePerformerRequest struct {
	Name        string  `json:"name"        binding:"omitempty,min=2,max=100"`
	Category    string  `json:"category"    binding:"omitempty"`
	Price       float64 `json:"price"       binding:"omitempty,gte=0"`
	Description string  `json:"description"`
}

// PerformerFilter holds query parameters for listing performers.
type PerformerFilter struct {
	Category string
	Page     int
	PageSize int
}
