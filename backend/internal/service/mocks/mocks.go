// Package mocks provides lightweight in-memory mock implementations of the
// domain repository interfaces, used exclusively in unit tests.
package mocks

import (
	"context"
	"sync"

	"epbms/internal/domain"
)

// ---- BookingRepository mock ----

// BookingRepo is a thread-safe in-memory mock of domain.BookingRepository.
type BookingRepo struct {
	mu       sync.RWMutex
	bookings map[uint]*domain.Booking
	nextID   uint
	// ConflictResult controls what FindConflicts returns (for deterministic tests).
	ConflictResult bool
	ConflictErr    error
}

func NewBookingRepo() *BookingRepo {
	return &BookingRepo{bookings: make(map[uint]*domain.Booking), nextID: 1}
}

func (m *BookingRepo) Create(ctx context.Context, b *domain.Booking) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b.ID = m.nextID
	m.nextID++
	cp := *b
	m.bookings[b.ID] = &cp
	return nil
}

func (m *BookingRepo) FindByID(ctx context.Context, id uint) (*domain.Booking, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.bookings[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *b
	return &cp, nil
}

func (m *BookingRepo) FindAll(ctx context.Context, filter domain.BookingFilter) ([]domain.Booking, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []domain.Booking
	for _, b := range m.bookings {
		if filter.PerformerID != 0 && b.PerformerID != filter.PerformerID {
			continue
		}
		if filter.ClientID != 0 && b.ClientID != filter.ClientID {
			continue
		}
		if filter.Status != "" && b.Status != filter.Status {
			continue
		}
		result = append(result, *b)
	}
	return result, int64(len(result)), nil
}

func (m *BookingRepo) UpdateStatus(ctx context.Context, id uint, status domain.BookingStatus, approvedBy *uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.bookings[id]
	if !ok {
		return domain.ErrNotFound
	}
	b.Status = status
	b.ApprovedBy = approvedBy
	return nil
}

func (m *BookingRepo) FindConflicts(ctx context.Context, performerID uint, eventDate, startTime, endTime string, excludeID uint) (bool, error) {
	return m.ConflictResult, m.ConflictErr
}

func (m *BookingRepo) Delete(ctx context.Context, id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.bookings[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.bookings, id)
	return nil
}

// ---- PerformerRepository mock ----

// PerformerRepo is a thread-safe in-memory mock of domain.PerformerRepository.
type PerformerRepo struct {
	mu         sync.RWMutex
	performers map[uint]*domain.Performer
	byUserID   map[uint]*domain.Performer
	nextID     uint
}

func NewPerformerRepo() *PerformerRepo {
	return &PerformerRepo{
		performers: make(map[uint]*domain.Performer),
		byUserID:   make(map[uint]*domain.Performer),
		nextID:     1,
	}
}

func (m *PerformerRepo) Create(ctx context.Context, p *domain.Performer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p.ID = m.nextID
	m.nextID++
	cp := *p
	m.performers[p.ID] = &cp
	m.byUserID[p.UserID] = &cp
	return nil
}

func (m *PerformerRepo) FindByID(ctx context.Context, id uint) (*domain.Performer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.performers[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (m *PerformerRepo) FindByUserID(ctx context.Context, userID uint) (*domain.Performer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.byUserID[userID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (m *PerformerRepo) FindAll(ctx context.Context, filter domain.PerformerFilter) ([]domain.Performer, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []domain.Performer
	for _, p := range m.performers {
		result = append(result, *p)
	}
	return result, int64(len(result)), nil
}

func (m *PerformerRepo) Update(ctx context.Context, p *domain.Performer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.performers[p.ID]; !ok {
		return domain.ErrNotFound
	}
	cp := *p
	m.performers[p.ID] = &cp
	return nil
}

func (m *PerformerRepo) Delete(ctx context.Context, id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.performers[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.performers, id)
	return nil
}

// ---- UserRepository mock ----

// UserRepo is a thread-safe in-memory mock of domain.UserRepository.
type UserRepo struct {
	mu     sync.RWMutex
	users  map[uint]*domain.User
	byEmail map[string]*domain.User
	nextID uint
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users:   make(map[uint]*domain.User),
		byEmail: make(map[string]*domain.User),
		nextID:  1,
	}
}

func (m *UserRepo) Create(ctx context.Context, u *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u.ID = m.nextID
	m.nextID++
	cp := *u
	m.users[u.ID] = &cp
	m.byEmail[u.Email] = &cp
	return nil
}

func (m *UserRepo) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (m *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.byEmail[email]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *u
	return &cp, nil
}
