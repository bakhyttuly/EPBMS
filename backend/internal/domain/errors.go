package domain

import "errors"

// Sentinel errors for domain-level error handling and mapping to HTTP status codes.
var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	ErrBookingConflict   = errors.New("booking time conflicts with an existing confirmed booking")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden: insufficient permissions")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidStatus     = errors.New("invalid status transition")
)
