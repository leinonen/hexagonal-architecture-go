package dto

import (
	"time"

	"github.com/leinonen/hexagonal-architecture-go/internal/domain"
)

type CreateUserDTO struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=1"`
}

type UpdateUserDTO struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
	Name  string `json:"name,omitempty" validate:"omitempty,min=1"`
}

type UserResponseDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToUserResponseDTO(user *domain.User) *UserResponseDTO {
	if user == nil {
		return nil
	}

	return &UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func ToUserResponseDTOs(users []*domain.User) []*UserResponseDTO {
	if users == nil {
		return nil
	}

	dtos := make([]*UserResponseDTO, len(users))
	for i, user := range users {
		dtos[i] = ToUserResponseDTO(user)
	}
	return dtos
}
