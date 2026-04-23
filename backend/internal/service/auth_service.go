package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"epbms/internal/domain"
	"epbms/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo      domain.UserRepository
	performerRepo domain.PerformerRepository
	log           *slog.Logger
}

// NewAuthService creates a new AuthService with the given dependencies.
func NewAuthService(userRepo domain.UserRepository, performerRepo domain.PerformerRepository, log *slog.Logger) domain.AuthService {
	return &authService{
		userRepo:      userRepo,
		performerRepo: performerRepo,
		log:           log,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.UserResponse, error) {
	// Check for duplicate email.
	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, domain.ErrConflict
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("authService.Register lookup: %w", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("authService.Register hash: %w", err)
	}

	user := &domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: string(hashed),
		Role:     req.Role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("authService.Register create user: %w", err)
	}

	// Automatically create a Performer profile when the role is performer.
	if req.Role == domain.RolePerformer {
		performer := &domain.Performer{
			UserID:      user.ID,
			Name:        user.FullName,
			Category:    "Not set",
			Price:       0,
			Description: "",
		}
		if err := s.performerRepo.Create(ctx, performer); err != nil {
			// Log but do not fail registration; the profile can be created later.
			s.log.Error("failed to auto-create performer profile", "user_id", user.ID, "error", err)
		}
	}

	s.log.Info("user registered", "user_id", user.ID, "role", user.Role)

	return &domain.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("authService.Login lookup: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("authService.Login generate token: %w", err)
	}

	s.log.Info("user logged in", "user_id", user.ID, "role", user.Role)

	return &domain.AuthResponse{
		Token: token,
		User: domain.UserResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}
