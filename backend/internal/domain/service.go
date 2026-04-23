package domain

import "context"

// AuthService defines the contract for authentication business logic.
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

// PerformerService defines the contract for performer business logic.
type PerformerService interface {
	GetAll(ctx context.Context, filter PerformerFilter) ([]Performer, int64, error)
	GetByID(ctx context.Context, id uint) (*Performer, error)
	Create(ctx context.Context, userID uint, req CreatePerformerRequest) (*Performer, error)
	Update(ctx context.Context, id uint, callerID uint, callerRole Role, req UpdatePerformerRequest) (*Performer, error)
	Delete(ctx context.Context, id uint) error
}

// BookingService defines the contract for booking business logic.
type BookingService interface {
	// CreateRequest is called by a CLIENT to create a pending booking.
	CreateRequest(ctx context.Context, clientID uint, req CreateBookingRequest) (*Booking, error)
	// UpdateStatus is called by ADMIN to approve/reject/complete a booking.
	UpdateStatus(ctx context.Context, bookingID uint, adminID uint, req UpdateBookingStatusRequest) (*Booking, error)
	// GetAll retrieves bookings filtered by the caller's role.
	GetAll(ctx context.Context, callerID uint, callerRole Role, filter BookingFilter) ([]Booking, int64, error)
	// GetByID retrieves a single booking, enforcing visibility rules.
	GetByID(ctx context.Context, bookingID uint, callerID uint, callerRole Role) (*Booking, error)
	// Delete removes a booking (admin only).
	Delete(ctx context.Context, bookingID uint) error
}
