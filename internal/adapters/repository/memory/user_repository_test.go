package memory_test

import (
	"context"
	"testing"

	"github.com/leinonen/hexagonal-architecture-go/internal/adapters/repository/memory"
	"github.com/leinonen/hexagonal-architecture-go/internal/domain"
	"github.com/leinonen/hexagonal-architecture-go/internal/errors"
)

func TestUserRepository_Create(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewUserRepository()

	user := &domain.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if user.ID == "" {
		t.Errorf("Create() should assign an ID to the user")
	}

	err = repo.Create(ctx, user)
	if !errors.IsConflict(err) {
		t.Errorf("Create() duplicate email should return conflict error")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewUserRepository()

	user := &domain.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	retrieved, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Errorf("GetByID() unexpected error: %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("GetByID() returned wrong user")
	}

	_, err = repo.GetByID(ctx, "non_existing")
	if !errors.IsNotFound(err) {
		t.Errorf("GetByID() non-existing user should return not found error")
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewUserRepository()

	user := &domain.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	retrieved, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Errorf("GetByEmail() unexpected error: %v", err)
	}

	if retrieved.Email != user.Email {
		t.Errorf("GetByEmail() returned wrong user")
	}

	_, err = repo.GetByEmail(ctx, "non@existing.com")
	if !errors.IsNotFound(err) {
		t.Errorf("GetByEmail() non-existing email should return not found error")
	}
}

func TestUserRepository_Update(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewUserRepository()

	user := &domain.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	user.Name = "Updated Name"
	err = repo.Update(ctx, user)
	if err != nil {
		t.Errorf("Update() unexpected error: %v", err)
	}

	retrieved, _ := repo.GetByID(ctx, user.ID)
	if retrieved.Name != "Updated Name" {
		t.Errorf("Update() did not update user name")
	}

	nonExisting := &domain.User{
		ID:    "non_existing",
		Email: "new@example.com",
		Name:  "New User",
	}
	err = repo.Update(ctx, nonExisting)
	if !errors.IsNotFound(err) {
		t.Errorf("Update() non-existing user should return not found error")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewUserRepository()

	user := &domain.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	err = repo.Delete(ctx, user.ID)
	if err != nil {
		t.Errorf("Delete() unexpected error: %v", err)
	}

	_, err = repo.GetByID(ctx, user.ID)
	if !errors.IsNotFound(err) {
		t.Errorf("Delete() user should not exist after deletion")
	}

	err = repo.Delete(ctx, "non_existing")
	if !errors.IsNotFound(err) {
		t.Errorf("Delete() non-existing user should return not found error")
	}
}

func TestUserRepository_List(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewUserRepository()

	for i := 0; i < 5; i++ {
		user := &domain.User{
			Email: string(rune('a'+i)) + "@example.com",
			Name:  "User " + string(rune('A'+i)),
		}
		repo.Create(ctx, user)
	}

	tests := []struct {
		name      string
		limit     int
		offset    int
		wantCount int
	}{
		{
			name:      "all users",
			limit:     10,
			offset:    0,
			wantCount: 5,
		},
		{
			name:      "limit 2",
			limit:     2,
			offset:    0,
			wantCount: 2,
		},
		{
			name:      "offset 2",
			limit:     10,
			offset:    2,
			wantCount: 3,
		},
		{
			name:      "limit 2 offset 3",
			limit:     2,
			offset:    3,
			wantCount: 2,
		},
		{
			name:      "offset beyond data",
			limit:     10,
			offset:    10,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := repo.List(ctx, tt.limit, tt.offset)
			if err != nil {
				t.Errorf("List() unexpected error: %v", err)
				return
			}

			if len(users) != tt.wantCount {
				t.Errorf("List() returned %d users, want %d", len(users), tt.wantCount)
			}
		})
	}
}

func TestUserRepository_Concurrency(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewUserRepository()

	done := make(chan bool)

	go func() {
		for i := 0; i < 100; i++ {
			user := &domain.User{
				Email: string(rune(i)) + "@example.com",
				Name:  "User",
			}
			repo.Create(ctx, user)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			repo.List(ctx, 10, 0)
		}
		done <- true
	}()

	<-done
	<-done

	users, err := repo.List(ctx, 200, 0)
	if err != nil {
		t.Errorf("List() after concurrent operations failed: %v", err)
	}

	if len(users) != 100 {
		t.Errorf("Expected 100 users after concurrent creates, got %d", len(users))
	}
}
