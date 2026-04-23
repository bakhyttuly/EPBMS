package repository

import (
	"context"
	"errors"
	"fmt"

	"epbms/internal/domain"
	"gorm.io/gorm"
)

type performerRepository struct {
	db *gorm.DB
}

// NewPerformerRepository creates a new instance of PerformerRepository backed by GORM.
func NewPerformerRepository(db *gorm.DB) domain.PerformerRepository {
	return &performerRepository{db: db}
}

func (r *performerRepository) Create(ctx context.Context, performer *domain.Performer) error {
	if err := r.db.WithContext(ctx).Create(performer).Error; err != nil {
		return fmt.Errorf("performerRepository.Create: %w", err)
	}
	return nil
}

func (r *performerRepository) FindByID(ctx context.Context, id uint) (*domain.Performer, error) {
	var performer domain.Performer
	err := r.db.WithContext(ctx).First(&performer, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("performerRepository.FindByID: %w", err)
	}
	return &performer, nil
}

func (r *performerRepository) FindByUserID(ctx context.Context, userID uint) (*domain.Performer, error) {
	var performer domain.Performer
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&performer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("performerRepository.FindByUserID: %w", err)
	}
	return &performer, nil
}

func (r *performerRepository) FindAll(ctx context.Context, filter domain.PerformerFilter) ([]domain.Performer, int64, error) {
	var performers []domain.Performer
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Performer{})
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("performerRepository.FindAll count: %w", err)
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

	if err := query.Offset(offset).Limit(pageSize).Find(&performers).Error; err != nil {
		return nil, 0, fmt.Errorf("performerRepository.FindAll: %w", err)
	}
	return performers, total, nil
}

func (r *performerRepository) Update(ctx context.Context, performer *domain.Performer) error {
	if err := r.db.WithContext(ctx).Save(performer).Error; err != nil {
		return fmt.Errorf("performerRepository.Update: %w", err)
	}
	return nil
}

func (r *performerRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Performer{}, id)
	if result.Error != nil {
		return fmt.Errorf("performerRepository.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
