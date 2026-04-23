package service

import (
	"context"
	"fmt"
	"log/slog"

	"epbms/internal/domain"
)

type performerService struct {
	performerRepo domain.PerformerRepository
	log           *slog.Logger
}

// NewPerformerService creates a new PerformerService with the given dependencies.
func NewPerformerService(performerRepo domain.PerformerRepository, log *slog.Logger) domain.PerformerService {
	return &performerService{
		performerRepo: performerRepo,
		log:           log,
	}
}

func (s *performerService) GetAll(ctx context.Context, filter domain.PerformerFilter) ([]domain.Performer, int64, error) {
	performers, total, err := s.performerRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("performerService.GetAll: %w", err)
	}
	return performers, total, nil
}

func (s *performerService) GetByID(ctx context.Context, id uint) (*domain.Performer, error) {
	performer, err := s.performerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("performerService.GetByID: %w", err)
	}
	return performer, nil
}

func (s *performerService) Create(ctx context.Context, userID uint, req domain.CreatePerformerRequest) (*domain.Performer, error) {
	performer := &domain.Performer{
		UserID:      userID,
		Name:        req.Name,
		Category:    req.Category,
		Price:       req.Price,
		Description: req.Description,
	}
	if err := s.performerRepo.Create(ctx, performer); err != nil {
		return nil, fmt.Errorf("performerService.Create: %w", err)
	}
	s.log.Info("performer profile created", "performer_id", performer.ID, "user_id", userID)
	return performer, nil
}

func (s *performerService) Update(ctx context.Context, id uint, callerID uint, callerRole domain.Role, req domain.UpdatePerformerRequest) (*domain.Performer, error) {
	performer, err := s.performerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("performerService.Update find: %w", err)
	}

	// A PERFORMER can only update their own profile.
	if callerRole == domain.RolePerformer && performer.UserID != callerID {
		return nil, domain.ErrForbidden
	}

	if req.Name != "" {
		performer.Name = req.Name
	}
	if req.Category != "" {
		performer.Category = req.Category
	}
	if req.Price >= 0 && req.Price != performer.Price {
		performer.Price = req.Price
	}
	if req.Description != "" {
		performer.Description = req.Description
	}

	if err := s.performerRepo.Update(ctx, performer); err != nil {
		return nil, fmt.Errorf("performerService.Update save: %w", err)
	}
	s.log.Info("performer profile updated", "performer_id", performer.ID)
	return performer, nil
}

func (s *performerService) Delete(ctx context.Context, id uint) error {
	if err := s.performerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("performerService.Delete: %w", err)
	}
	s.log.Info("performer profile deleted", "performer_id", id)
	return nil
}
