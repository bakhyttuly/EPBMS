package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	FullName string `gorm:"not null" json:"full_name"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"not null" json:"role"`
}
