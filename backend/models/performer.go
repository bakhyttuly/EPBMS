package models

type Performer struct {
	ID                 uint    `gorm:"primaryKey" json:"id"`
	UserID             *uint   `gorm:"uniqueIndex" json:"user_id,omitempty"`
	Name               string  `gorm:"not null" json:"name"`
	Category           string  `gorm:"not null" json:"category"`
	Price              float64 `gorm:"not null" json:"price"`
	Description        string  `json:"description"`
	AvailabilityStatus string  `json:"availability_status"`
}
