package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/leinonen/hexagonal-architecture-go/internal/application"
	"github.com/leinonen/hexagonal-architecture-go/internal/domain"
	"github.com/leinonen/hexagonal-architecture-go/internal/errors"
)

type mockUserRepository struct {
	users         map[string]*domain.User
	createErr     error
	getByIDErr    error
	getByEmailErr error
	updateErr     error
	deleteErr     error
	listErr       error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	user.ID = "test_id_1"
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	user, exists := m.users[id]
	if !exists {
		return nil, errors.NewNotFoundError("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailErr != nil {
		return nil, m.getByEmailErr
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.NewNotFoundError("user not found")
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.users, id)
	return nil
}

func (m *mockUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		email     string
		userName  string
		setupMock func(*mockUserRepository)
		wantErr   bool
		errType   errors.ErrorType
	}{
		{
			name:      "successful creation",
			email:     "test@example.com",
			userName:  "Test User",
			setupMock: func(m *mockUserRepository) {},
			wantErr:   false,
		},
		{
			name:     "duplicate email",
			email:    "existing@example.com",
			userName: "New User",
			setupMock: func(m *mockUserRepository) {
				m.users["existing_id"] = &domain.User{
					ID:    "existing_id",
					Email: "existing@example.com",
					Name:  "Existing User",
				}
			},
			wantErr: true,
			errType: errors.Conflict,
		},
		{
			name:      "invalid email",
			email:     "",
			userName:  "Test User",
			setupMock: func(m *mockUserRepository) {},
			wantErr:   true,
			errType:   errors.Validation,
		},
		{
			name:      "invalid name",
			email:     "test@example.com",
			userName:  "",
			setupMock: func(m *mockUserRepository) {},
			wantErr:   true,
			errType:   errors.Validation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			tt.setupMock(repo)

			service := application.NewUserService(repo)
			user, err := service.CreateUser(ctx, tt.email, tt.userName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateUser() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateUser() unexpected error: %v", err)
				return
			}

			if user.Email != tt.email {
				t.Errorf("CreateUser() email = %v, want %v", user.Email, tt.email)
			}

			if user.Name != tt.userName {
				t.Errorf("CreateUser() name = %v, want %v", user.Name, tt.userName)
			}

			if user.ID == "" {
				t.Errorf("CreateUser() ID should not be empty")
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		setupMock func(*mockUserRepository)
		wantErr   bool
		wantEmail string
	}{
		{
			name:   "existing user",
			userID: "test_id_1",
			setupMock: func(m *mockUserRepository) {
				m.users["test_id_1"] = &domain.User{
					ID:        "test_id_1",
					Email:     "test@example.com",
					Name:      "Test User",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
			},
			wantErr:   false,
			wantEmail: "test@example.com",
		},
		{
			name:      "non-existing user",
			userID:    "non_existing",
			setupMock: func(m *mockUserRepository) {},
			wantErr:   true,
		},
		{
			name:      "empty ID",
			userID:    "",
			setupMock: func(m *mockUserRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			tt.setupMock(repo)

			service := application.NewUserService(repo)
			user, err := service.GetUser(ctx, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUser() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GetUser() unexpected error: %v", err)
				return
			}

			if user.Email != tt.wantEmail {
				t.Errorf("GetUser() email = %v, want %v", user.Email, tt.wantEmail)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		newEmail  string
		newName   string
		setupMock func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name:     "successful update name",
			userID:   "test_id_1",
			newEmail: "",
			newName:  "Updated Name",
			setupMock: func(m *mockUserRepository) {
				m.users["test_id_1"] = &domain.User{
					ID:    "test_id_1",
					Email: "test@example.com",
					Name:  "Original Name",
				}
			},
			wantErr: false,
		},
		{
			name:     "successful update email",
			userID:   "test_id_1",
			newEmail: "newemail@example.com",
			newName:  "",
			setupMock: func(m *mockUserRepository) {
				m.users["test_id_1"] = &domain.User{
					ID:    "test_id_1",
					Email: "test@example.com",
					Name:  "Test User",
				}
			},
			wantErr: false,
		},
		{
			name:     "email already in use",
			userID:   "test_id_1",
			newEmail: "existing@example.com",
			newName:  "",
			setupMock: func(m *mockUserRepository) {
				m.users["test_id_1"] = &domain.User{
					ID:    "test_id_1",
					Email: "test@example.com",
					Name:  "Test User",
				}
				m.users["test_id_2"] = &domain.User{
					ID:    "test_id_2",
					Email: "existing@example.com",
					Name:  "Existing User",
				}
			},
			wantErr: true,
		},
		{
			name:      "non-existing user",
			userID:    "non_existing",
			newEmail:  "new@example.com",
			newName:   "New Name",
			setupMock: func(m *mockUserRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			tt.setupMock(repo)

			service := application.NewUserService(repo)
			user, err := service.UpdateUser(ctx, tt.userID, tt.newEmail, tt.newName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateUser() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateUser() unexpected error: %v", err)
				return
			}

			if tt.newEmail != "" && user.Email != tt.newEmail {
				t.Errorf("UpdateUser() email = %v, want %v", user.Email, tt.newEmail)
			}

			if tt.newName != "" && user.Name != tt.newName {
				t.Errorf("UpdateUser() name = %v, want %v", user.Name, tt.newName)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		setupMock func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name:   "successful deletion",
			userID: "test_id_1",
			setupMock: func(m *mockUserRepository) {
				m.users["test_id_1"] = &domain.User{
					ID:    "test_id_1",
					Email: "test@example.com",
					Name:  "Test User",
				}
			},
			wantErr: false,
		},
		{
			name:      "empty ID",
			userID:    "",
			setupMock: func(m *mockUserRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			tt.setupMock(repo)

			service := application.NewUserService(repo)
			err := service.DeleteUser(ctx, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteUser() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteUser() unexpected error: %v", err)
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	ctx := context.Background()

	repo := newMockUserRepository()
	repo.users["1"] = &domain.User{ID: "1", Email: "user1@example.com"}
	repo.users["2"] = &domain.User{ID: "2", Email: "user2@example.com"}
	repo.users["3"] = &domain.User{ID: "3", Email: "user3@example.com"}

	service := application.NewUserService(repo)

	tests := []struct {
		name      string
		limit     int
		offset    int
		wantCount int
	}{
		{
			name:      "default limit",
			limit:     0,
			offset:    0,
			wantCount: 3,
		},
		{
			name:      "custom limit",
			limit:     2,
			offset:    0,
			wantCount: 3,
		},
		{
			name:      "negative offset becomes zero",
			limit:     10,
			offset:    -5,
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := service.ListUsers(ctx, tt.limit, tt.offset)
			if err != nil {
				t.Errorf("ListUsers() unexpected error: %v", err)
				return
			}

			if len(users) != tt.wantCount {
				t.Errorf("ListUsers() returned %d users, want %d", len(users), tt.wantCount)
			}
		})
	}
}
