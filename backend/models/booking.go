package models

type Booking struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	OrganizerID uint   `gorm:"not null" json:"organizer_id"`
	PerformerID uint   `gorm:"not null" json:"performer_id"`
	ClientName  string `gorm:"not null" json:"client_name"`
	EventDate   string `gorm:"not null" json:"event_date"`
	StartTime   string `gorm:"not null" json:"start_time"`
	EndTime     string `gorm:"not null" json:"end_time"`
	Status      string `json:"status"`

	Performer Performer `gorm:"foreignKey:PerformerID" json:"performer"`
	Organizer User      `gorm:"foreignKey:OrganizerID" json:"-"`
}
