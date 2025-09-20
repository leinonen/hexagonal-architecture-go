package application

import (
	"context"
	"time"

	"github.com/leinonen/hexagonal-architecture-go/internal/domain"
	"github.com/leinonen/hexagonal-architecture-go/internal/errors"
	"github.com/leinonen/hexagonal-architecture-go/internal/ports"
)

type UserService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
	user, err := domain.NewUser(email, name)
	if err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.NewConflictError("user with this email already exists")
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.NewValidationError("user ID is required")
	}

	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "" {
		return nil, errors.NewValidationError("email is required")
	}

	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, id, email, name string) (*domain.User, error) {
	if id == "" {
		return nil, errors.NewValidationError("user ID is required")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if email != "" && email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, email)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.NewConflictError("email already in use")
		}
		user.Email = email
	}

	if name != "" {
		user.Name = name
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.NewValidationError("user ID is required")
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.userRepo.List(ctx, limit, offset)
}
