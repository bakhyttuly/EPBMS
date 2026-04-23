package service_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"epbms/internal/domain"
	"epbms/internal/service"
	"epbms/internal/service/mocks"
)

// newTestLogger returns a discard logger so test output is clean.
func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

// seedPerformer inserts a performer into the mock repo and returns it.
func seedPerformer(repo *mocks.PerformerRepo, userID uint) *domain.Performer {
	p := &domain.Performer{UserID: userID, Name: "Test Performer", Category: "Music", Price: 100}
	_ = repo.Create(context.Background(), p)
	return p
}

// ============================================================
// Booking Conflict Detection Tests
// ============================================================

// TestCreateRequest_NoConflict verifies that a booking is created successfully
// when the mock repository reports no conflict.
func TestCreateRequest_NoConflict(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()
	performer := seedPerformer(performerRepo, 10)

	bookingRepo.ConflictResult = false

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	req := domain.CreateBookingRequest{
		PerformerID: performer.ID,
		EventDate:   "2025-12-01",
		StartTime:   "10:00",
		EndTime:     "12:00",
	}

	booking, err := svc.CreateRequest(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if booking.Status != domain.StatusPending {
		t.Errorf("expected status %q, got %q", domain.StatusPending, booking.Status)
	}
	if booking.ClientID != 1 {
		t.Errorf("expected client_id 1, got %d", booking.ClientID)
	}
}

// TestCreateRequest_PerformerNotFound verifies that creating a booking for a
// non-existent performer returns ErrNotFound.
func TestCreateRequest_PerformerNotFound(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	req := domain.CreateBookingRequest{
		PerformerID: 999, // does not exist
		EventDate:   "2025-12-01",
		StartTime:   "10:00",
		EndTime:     "12:00",
	}

	_, err := svc.CreateRequest(context.Background(), 1, req)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

// TestUpdateStatus_ConfirmWithConflict verifies that confirming a booking
// returns ErrBookingConflict when the repository detects an overlap.
func TestUpdateStatus_ConfirmWithConflict(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()
	performer := seedPerformer(performerRepo, 10)

	// Create a pending booking to be confirmed.
	pending := &domain.Booking{
		PerformerID: performer.ID,
		ClientID:    1,
		EventDate:   "2025-12-01",
		StartTime:   "10:00",
		EndTime:     "12:00",
		Status:      domain.StatusPending,
	}
	_ = bookingRepo.Create(context.Background(), pending)

	// Simulate a conflict detected by the repository.
	bookingRepo.ConflictResult = true

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	req := domain.UpdateBookingStatusRequest{Status: domain.StatusConfirmed}
	_, err := svc.UpdateStatus(context.Background(), pending.ID, 99, req)

	if !errors.Is(err, domain.ErrBookingConflict) {
		t.Errorf("expected ErrBookingConflict, got: %v", err)
	}
}

// TestUpdateStatus_ConfirmNoConflict verifies that a booking transitions to
// "confirmed" when no conflict is detected.
func TestUpdateStatus_ConfirmNoConflict(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()
	performer := seedPerformer(performerRepo, 10)

	pending := &domain.Booking{
		PerformerID: performer.ID,
		ClientID:    1,
		EventDate:   "2025-12-01",
		StartTime:   "10:00",
		EndTime:     "12:00",
		Status:      domain.StatusPending,
	}
	_ = bookingRepo.Create(context.Background(), pending)

	bookingRepo.ConflictResult = false

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	req := domain.UpdateBookingStatusRequest{Status: domain.StatusConfirmed}
	updated, err := svc.UpdateStatus(context.Background(), pending.ID, 99, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if updated.Status != domain.StatusConfirmed {
		t.Errorf("expected status %q, got %q", domain.StatusConfirmed, updated.Status)
	}
}

// TestUpdateStatus_RejectBooking verifies that a booking can be rejected
// without a conflict check.
func TestUpdateStatus_RejectBooking(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()
	performer := seedPerformer(performerRepo, 10)

	pending := &domain.Booking{
		PerformerID: performer.ID,
		ClientID:    1,
		EventDate:   "2025-12-01",
		StartTime:   "10:00",
		EndTime:     "12:00",
		Status:      domain.StatusPending,
	}
	_ = bookingRepo.Create(context.Background(), pending)

	// Even if conflict is true, rejection should not check for conflicts.
	bookingRepo.ConflictResult = true

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	req := domain.UpdateBookingStatusRequest{Status: domain.StatusRejected}
	updated, err := svc.UpdateStatus(context.Background(), pending.ID, 99, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if updated.Status != domain.StatusRejected {
		t.Errorf("expected status %q, got %q", domain.StatusRejected, updated.Status)
	}
}

// TestUpdateStatus_InvalidTransition verifies that transitioning from a
// terminal state (completed/rejected) returns ErrInvalidStatus.
func TestUpdateStatus_InvalidTransition(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()
	performer := seedPerformer(performerRepo, 10)

	completed := &domain.Booking{
		PerformerID: performer.ID,
		ClientID:    1,
		EventDate:   "2025-12-01",
		StartTime:   "10:00",
		EndTime:     "12:00",
		Status:      domain.StatusCompleted,
	}
	_ = bookingRepo.Create(context.Background(), completed)

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	req := domain.UpdateBookingStatusRequest{Status: domain.StatusConfirmed}
	_, err := svc.UpdateStatus(context.Background(), completed.ID, 99, req)

	if !errors.Is(err, domain.ErrInvalidStatus) {
		t.Errorf("expected ErrInvalidStatus, got: %v", err)
	}
}

// ============================================================
// Role-Based Visibility Tests
// ============================================================

// TestGetAll_ClientSeesOnlyOwnBookings verifies that a CLIENT can only retrieve
// bookings where they are the client.
func TestGetAll_ClientSeesOnlyOwnBookings(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()
	performer := seedPerformer(performerRepo, 10)

	// Booking owned by client 1.
	_ = bookingRepo.Create(context.Background(), &domain.Booking{
		PerformerID: performer.ID, ClientID: 1, Status: domain.StatusPending,
		EventDate: "2025-12-01", StartTime: "10:00", EndTime: "12:00",
	})
	// Booking owned by client 2.
	_ = bookingRepo.Create(context.Background(), &domain.Booking{
		PerformerID: performer.ID, ClientID: 2, Status: domain.StatusPending,
		EventDate: "2025-12-02", StartTime: "10:00", EndTime: "12:00",
	})

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	bookings, total, err := svc.GetAll(context.Background(), 1, domain.RoleClient, domain.BookingFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 booking for client 1, got %d", total)
	}
	if bookings[0].ClientID != 1 {
		t.Errorf("expected booking for client 1, got client %d", bookings[0].ClientID)
	}
}

// TestGetAll_PerformerSeesOnlyConfirmedOwnBookings verifies that a PERFORMER
// can only see CONFIRMED bookings assigned to them.
func TestGetAll_PerformerSeesOnlyConfirmedOwnBookings(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()

	// Two performers: user 10 -> performer 1, user 20 -> performer 2.
	p1 := seedPerformer(performerRepo, 10)
	p2 := seedPerformer(performerRepo, 20)

	// Confirmed booking for performer 1.
	_ = bookingRepo.Create(context.Background(), &domain.Booking{
		PerformerID: p1.ID, ClientID: 1, Status: domain.StatusConfirmed,
		EventDate: "2025-12-01", StartTime: "10:00", EndTime: "12:00",
	})
	// Pending booking for performer 1 — should NOT be visible to performer.
	_ = bookingRepo.Create(context.Background(), &domain.Booking{
		PerformerID: p1.ID, ClientID: 2, Status: domain.StatusPending,
		EventDate: "2025-12-02", StartTime: "10:00", EndTime: "12:00",
	})
	// Confirmed booking for performer 2 — should NOT be visible to performer 1.
	_ = bookingRepo.Create(context.Background(), &domain.Booking{
		PerformerID: p2.ID, ClientID: 3, Status: domain.StatusConfirmed,
		EventDate: "2025-12-03", StartTime: "10:00", EndTime: "12:00",
	})

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	// Caller is user 10 (performer 1).
	bookings, total, err := svc.GetAll(context.Background(), 10, domain.RolePerformer, domain.BookingFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 confirmed booking for performer 1, got %d", total)
	}
	if bookings[0].PerformerID != p1.ID {
		t.Errorf("expected performer_id %d, got %d", p1.ID, bookings[0].PerformerID)
	}
	if bookings[0].Status != domain.StatusConfirmed {
		t.Errorf("expected status %q, got %q", domain.StatusConfirmed, bookings[0].Status)
	}
}

// TestGetAll_AdminSeesAllBookings verifies that an ADMIN can see all bookings
// regardless of status or ownership.
func TestGetAll_AdminSeesAllBookings(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()
	p := seedPerformer(performerRepo, 10)

	statuses := []domain.BookingStatus{
		domain.StatusPending, domain.StatusConfirmed, domain.StatusRejected, domain.StatusCompleted,
	}
	for i, s := range statuses {
		_ = bookingRepo.Create(context.Background(), &domain.Booking{
			PerformerID: p.ID, ClientID: uint(i + 1), Status: s,
			EventDate: "2025-12-01", StartTime: "10:00", EndTime: "11:00",
		})
	}

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	_, total, err := svc.GetAll(context.Background(), 99, domain.RoleAdmin, domain.BookingFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != int64(len(statuses)) {
		t.Errorf("expected %d bookings for admin, got %d", len(statuses), total)
	}
}

// TestGetByID_ClientForbiddenOnOtherBooking verifies that a CLIENT cannot
// access a booking that belongs to a different client.
func TestGetByID_ClientForbiddenOnOtherBooking(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()
	p := seedPerformer(performerRepo, 10)

	booking := &domain.Booking{
		PerformerID: p.ID, ClientID: 2, Status: domain.StatusPending,
		EventDate: "2025-12-01", StartTime: "10:00", EndTime: "12:00",
	}
	_ = bookingRepo.Create(context.Background(), booking)

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	// Caller is client 1, but booking belongs to client 2.
	_, err := svc.GetByID(context.Background(), booking.ID, 1, domain.RoleClient)
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got: %v", err)
	}
}

// TestGetByID_PerformerForbiddenOnPendingBooking verifies that a PERFORMER
// cannot access a booking that is still pending (not yet confirmed).
func TestGetByID_PerformerForbiddenOnPendingBooking(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()
	p := seedPerformer(performerRepo, 10)

	booking := &domain.Booking{
		PerformerID: p.ID, ClientID: 1, Status: domain.StatusPending,
		EventDate: "2025-12-01", StartTime: "10:00", EndTime: "12:00",
	}
	_ = bookingRepo.Create(context.Background(), booking)

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	// Caller is user 10 (the performer), but booking is pending.
	_, err := svc.GetByID(context.Background(), booking.ID, 10, domain.RolePerformer)
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got: %v", err)
	}
}

// TestGetByID_PerformerCanAccessConfirmedOwnBooking verifies that a PERFORMER
// can access a CONFIRMED booking assigned to them.
func TestGetByID_PerformerCanAccessConfirmedOwnBooking(t *testing.T) {
	bookingRepo := mocks.NewBookingRepo()
	performerRepo := mocks.NewPerformerRepo()
	p := seedPerformer(performerRepo, 10)

	booking := &domain.Booking{
		PerformerID: p.ID, ClientID: 1, Status: domain.StatusConfirmed,
		EventDate: "2025-12-01", StartTime: "10:00", EndTime: "12:00",
	}
	_ = bookingRepo.Create(context.Background(), booking)

	svc := service.NewBookingService(bookingRepo, performerRepo, newTestLogger())

	result, err := svc.GetByID(context.Background(), booking.ID, 10, domain.RolePerformer)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.ID != booking.ID {
		t.Errorf("expected booking ID %d, got %d", booking.ID, result.ID)
	}
}
