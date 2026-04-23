package domain

import "time"

// Role represents the access level of a user in the system.
type Role string

const (
	RoleAdmin     Role = "admin"
	RolePerformer Role = "performer"
	RoleClient    Role = "client"
)

// User is the core user entity.
type User struct {
	ID        uint      `gorm:"primaryKey"                  json:"id"`
	FullName  string    `gorm:"not null"                    json:"full_name"`
	Email     string    `gorm:"uniqueIndex;not null"        json:"email"`
	Password  string    `gorm:"not null"                    json:"-"`
	Role      Role      `gorm:"type:varchar(20);not null"   json:"role"`
	CreatedAt time.Time `                                   json:"created_at"`
	UpdatedAt time.Time `                                   json:"updated_at"`
}

// RegisterRequest is the payload for user registration.
type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8"`
	Role     Role   `json:"role"      binding:"required,oneof=admin performer client"`
}

// LoginRequest is the payload for user login.
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is returned after a successful login.
type AuthResponse struct {
	Token string `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse is the safe public representation of a user (no password).
type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
}
