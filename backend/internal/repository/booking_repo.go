package repository

import (
	"context"
	"errors"
	"fmt"

	"epbms/internal/domain"
	"gorm.io/gorm"
)

type bookingRepository struct {
	db *gorm.DB
}

// NewBookingRepository creates a new instance of BookingRepository backed by GORM.
func NewBookingRepository(db *gorm.DB) domain.BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	if err := r.db.WithContext(ctx).Create(booking).Error; err != nil {
		return fmt.Errorf("bookingRepository.Create: %w", err)
	}
	return nil
}

func (r *bookingRepository) FindByID(ctx context.Context, id uint) (*domain.Booking, error) {
	var booking domain.Booking
	err := r.db.WithContext(ctx).
		Preload("Performer").
		Preload("Client").
		First(&booking, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("bookingRepository.FindByID: %w", err)
	}
	return &booking, nil
}

func (r *bookingRepository) FindAll(ctx context.Context, filter domain.BookingFilter) ([]domain.Booking, int64, error) {
	var bookings []domain.Booking
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Booking{}).
		Preload("Performer")

	if filter.PerformerID != 0 {
		query = query.Where("performer_id = ?", filter.PerformerID)
	}
	if filter.ClientID != 0 {
		query = query.Where("client_id = ?", filter.ClientID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.EventDate != "" {
		query = query.Where("event_date = ?", filter.EventDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("bookingRepository.FindAll count: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := query.Order("event_date asc, start_time asc").Offset(offset).Limit(pageSize).Find(&bookings).Error; err != nil {
		return nil, 0, fmt.Errorf("bookingRepository.FindAll: %w", err)
	}
	return bookings, total, nil
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id uint, status domain.BookingStatus, approvedBy *uint) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if approvedBy != nil {
		updates["approved_by"] = approvedBy
		updates["approved_at"] = gorm.Expr("NOW()")
	}
	result := r.db.WithContext(ctx).Model(&domain.Booking{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("bookingRepository.UpdateStatus: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// FindConflicts checks if a confirmed booking exists for the same performer that overlaps
// the requested time window. The overlap condition is:
//
//	(start_time < existing_end_time) AND (end_time > existing_start_time)
//
// This is the standard interval overlap check. An optional excludeID allows
// skipping a specific booking (used when updating an existing booking).
func (r *bookingRepository) FindConflicts(ctx context.Context, performerID uint, eventDate, startTime, endTime string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.Booking{}).
		Where("performer_id = ?", performerID).
		Where("event_date = ?", eventDate).
		Where("status IN ?", []domain.BookingStatus{domain.StatusConfirmed, domain.StatusPending}).
		Where("start_time < ? AND end_time > ?", endTime, startTime)

	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("bookingRepository.FindConflicts: %w", err)
	}
	return count > 0, nil
}

func (r *bookingRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Booking{}, id)
	if result.Error != nil {
		return fmt.Errorf("bookingRepository.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
