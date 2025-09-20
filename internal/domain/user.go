package domain

import (
	"errors"
	"time"
)

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, name string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	return &User{
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}