package domain_test

import (
	"testing"
	"time"

	"github.com/leinonen/hexagonal-architecture-go/internal/domain"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		userName string
		wantErr  bool
	}{
		{
			name:     "valid user",
			email:    "test@example.com",
			userName: "Test User",
			wantErr:  false,
		},
		{
			name:     "empty email",
			email:    "",
			userName: "Test User",
			wantErr:  true,
		},
		{
			name:     "empty name",
			email:    "test@example.com",
			userName: "",
			wantErr:  true,
		},
		{
			name:     "both empty",
			email:    "",
			userName: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := domain.NewUser(tt.email, tt.userName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error: %v", err)
				return
			}

			if user.Email != tt.email {
				t.Errorf("NewUser() email = %v, want %v", user.Email, tt.email)
			}

			if user.Name != tt.userName {
				t.Errorf("NewUser() name = %v, want %v", user.Name, tt.userName)
			}

			if user.CreatedAt.IsZero() {
				t.Errorf("NewUser() CreatedAt should not be zero")
			}

			if user.UpdatedAt.IsZero() {
				t.Errorf("NewUser() UpdatedAt should not be zero")
			}

			if user.CreatedAt.After(time.Now()) {
				t.Errorf("NewUser() CreatedAt should not be in the future")
			}
		})
	}
}
