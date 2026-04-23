package domain

import "context"

// UserRepository defines the contract for user data persistence.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// PerformerRepository defines the contract for performer data persistence.
type PerformerRepository interface {
	Create(ctx context.Context, performer *Performer) error
	FindByID(ctx context.Context, id uint) (*Performer, error)
	FindByUserID(ctx context.Context, userID uint) (*Performer, error)
	FindAll(ctx context.Context, filter PerformerFilter) ([]Performer, int64, error)
	Update(ctx context.Context, performer *Performer) error
	Delete(ctx context.Context, id uint) error
}

// BookingRepository defines the contract for booking data persistence.
type BookingRepository interface {
	Create(ctx context.Context, booking *Booking) error
	FindByID(ctx context.Context, id uint) (*Booking, error)
	FindAll(ctx context.Context, filter BookingFilter) ([]Booking, int64, error)
	UpdateStatus(ctx context.Context, id uint, status BookingStatus, approvedBy *uint) error
	FindConflicts(ctx context.Context, performerID uint, eventDate, startTime, endTime string, excludeID uint) (bool, error)
	Delete(ctx context.Context, id uint) error
}
