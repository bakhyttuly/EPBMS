package service

import (
	"context"
	"fmt"
	"log/slog"

	"epbms/internal/domain"
)

type bookingService struct {
	bookingRepo   domain.BookingRepository
	performerRepo domain.PerformerRepository
	log           *slog.Logger
}

// NewBookingService creates a new BookingService with the given dependencies.
func NewBookingService(bookingRepo domain.BookingRepository, performerRepo domain.PerformerRepository, log *slog.Logger) domain.BookingService {
	return &bookingService{
		bookingRepo:   bookingRepo,
		performerRepo: performerRepo,
		log:           log,
	}
}

// CreateRequest is called by a CLIENT to submit a new booking request.
// The booking is created with status "pending" and does NOT yet check for conflicts,
// as conflicts are only relevant once the booking is confirmed by an ADMIN.
func (s *bookingService) CreateRequest(ctx context.Context, clientID uint, req domain.CreateBookingRequest) (*domain.Booking, error) {
	// Validate that the performer exists.
	_, err := s.performerRepo.FindByID(ctx, req.PerformerID)
	if err != nil {
		return nil, fmt.Errorf("bookingService.CreateRequest find performer: %w", err)
	}

	booking := &domain.Booking{
		PerformerID: req.PerformerID,
		ClientID:    clientID,
		EventDate:   req.EventDate,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      domain.StatusPending,
		Notes:       req.Notes,
	}

	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, fmt.Errorf("bookingService.CreateRequest create: %w", err)
	}

	s.log.Info("booking request created", "booking_id", booking.ID, "client_id", clientID, "performer_id", req.PerformerID)
	return booking, nil
}

// UpdateStatus is called by ADMIN to transition a booking's status.
// When confirming a booking, conflict detection is performed to ensure no
// overlapping confirmed/pending bookings exist for the same performer.
func (s *bookingService) UpdateStatus(ctx context.Context, bookingID uint, adminID uint, req domain.UpdateBookingStatusRequest) (*domain.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("bookingService.UpdateStatus find: %w", err)
	}

	// Guard against invalid state transitions.
	if booking.Status == domain.StatusCompleted || booking.Status == domain.StatusRejected {
		return nil, fmt.Errorf("%w: cannot change status from %s", domain.ErrInvalidStatus, booking.Status)
	}

	// Perform conflict detection only when confirming a booking.
	if req.Status == domain.StatusConfirmed {
		conflict, err := s.bookingRepo.FindConflicts(
			ctx,
			booking.PerformerID,
			booking.EventDate,
			booking.StartTime,
			booking.EndTime,
			bookingID, // exclude the current booking from the conflict check
		)
		if err != nil {
			return nil, fmt.Errorf("bookingService.UpdateStatus conflict check: %w", err)
		}
		if conflict {
			return nil, domain.ErrBookingConflict
		}
	}

	var approverID *uint
	if req.Status == domain.StatusConfirmed || req.Status == domain.StatusRejected {
		approverID = &adminID
	}

	if err := s.bookingRepo.UpdateStatus(ctx, bookingID, req.Status, approverID); err != nil {
		return nil, fmt.Errorf("bookingService.UpdateStatus update: %w", err)
	}

	s.log.Info("booking status updated",
		"booking_id", bookingID,
		"new_status", req.Status,
		"admin_id", adminID,
	)

	// Return the updated booking with preloaded associations.
	updated, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("bookingService.UpdateStatus reload: %w", err)
	}
	return updated, nil
}

// GetAll retrieves bookings with role-based visibility enforcement.
//
//   - ADMIN: sees all bookings, with optional filters.
//   - CLIENT: sees only their own bookings (any status).
//   - PERFORMER: sees only CONFIRMED bookings assigned to them.
func (s *bookingService) GetAll(ctx context.Context, callerID uint, callerRole domain.Role, filter domain.BookingFilter) ([]domain.Booking, int64, error) {
	switch callerRole {
	case domain.RoleClient:
		filter.ClientID = callerID
	case domain.RolePerformer:
		// Find the performer profile linked to this user.
		performer, err := s.performerRepo.FindByUserID(ctx, callerID)
		if err != nil {
			return nil, 0, fmt.Errorf("bookingService.GetAll find performer: %w", err)
		}
		filter.PerformerID = performer.ID
		filter.Status = domain.StatusConfirmed
	// ADMIN: no additional filter constraints; use whatever was passed in.
	}

	bookings, total, err := s.bookingRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("bookingService.GetAll: %w", err)
	}
	return bookings, total, nil
}

// GetByID retrieves a single booking, enforcing visibility rules per role.
func (s *bookingService) GetByID(ctx context.Context, bookingID uint, callerID uint, callerRole domain.Role) (*domain.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("bookingService.GetByID: %w", err)
	}

	switch callerRole {
	case domain.RoleClient:
		if booking.ClientID != callerID {
			return nil, domain.ErrForbidden
		}
	case domain.RolePerformer:
		performer, err := s.performerRepo.FindByUserID(ctx, callerID)
		if err != nil {
			return nil, fmt.Errorf("bookingService.GetByID find performer: %w", err)
		}
		if booking.PerformerID != performer.ID || booking.Status != domain.StatusConfirmed {
			return nil, domain.ErrForbidden
		}
	// ADMIN: full access, no additional checks.
	}

	return booking, nil
}

// Delete removes a booking record. Admin-only; enforced at the handler level.
func (s *bookingService) Delete(ctx context.Context, bookingID uint) error {
	if err := s.bookingRepo.Delete(ctx, bookingID); err != nil {
		return fmt.Errorf("bookingService.Delete: %w", err)
	}
	s.log.Info("booking deleted", "booking_id", bookingID)
	return nil
}
