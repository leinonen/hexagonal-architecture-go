package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpHandler "github.com/leinonen/hexagonal-architecture-go/internal/adapters/http"
	"github.com/leinonen/hexagonal-architecture-go/internal/application"
	"github.com/leinonen/hexagonal-architecture-go/internal/domain"
	"github.com/leinonen/hexagonal-architecture-go/internal/errors"
)

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]*domain.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		user.ID = "test_id"
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.NewNotFoundError("user not found")
	}
	return user, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.NewNotFoundError("user not found")
}

func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return errors.NewNotFoundError("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	if _, exists := m.users[id]; !exists {
		return errors.NewNotFoundError("user not found")
	}
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

type mockWeatherService struct {
	weather map[string]*domain.Weather
}

func newMockWeatherService() *mockWeatherService {
	return &mockWeatherService{
		weather: map[string]*domain.Weather{
			"London": {
				City:        "London",
				Temperature: 15.5,
				Description: "Cloudy",
				Humidity:    70,
				WindSpeed:   5.5,
			},
		},
	}
}

func (m *mockWeatherService) GetWeather(ctx context.Context, city string) (*domain.Weather, error) {
	weather, exists := m.weather[city]
	if !exists {
		return nil, errors.NewNotFoundError("city not found")
	}
	return weather, nil
}

func TestHandler_CreateUser(t *testing.T) {
	userRepo := newMockUserRepo()
	userService := application.NewUserService(userRepo)
	weatherService := application.NewWeatherService(newMockWeatherService(), userRepo)
	handler := httpHandler.NewHandler(userService, weatherService)

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid user",
			body: map[string]string{
				"email": "test@example.com",
				"name":  "Test User",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing email",
			body: map[string]string{
				"name": "Test User",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing name",
			body: map[string]string{
				"email": "test@example.com",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.CreateUser(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateUser() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusCreated {
				var response map[string]interface{}
				json.NewDecoder(w.Body).Decode(&response)

				if response["id"] == "" {
					t.Errorf("CreateUser() response missing ID")
				}

				if response["email"] != tt.body["email"] {
					t.Errorf("CreateUser() email = %v, want %v", response["email"], tt.body["email"])
				}
			}
		})
	}
}

func TestHandler_GetUser(t *testing.T) {
	userRepo := newMockUserRepo()
	testUser, _ := domain.NewUser("test@example.com", "Test User")
	testUser.ID = "test_id"
	userRepo.users["test_id"] = testUser

	userService := application.NewUserService(userRepo)
	weatherService := application.NewWeatherService(newMockWeatherService(), userRepo)
	handler := httpHandler.NewHandler(userService, weatherService)

	tests := []struct {
		name       string
		userID     string
		wantStatus int
	}{
		{
			name:       "existing user",
			userID:     "test_id",
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing user",
			userID:     "non_existing",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/users/"+tt.userID, nil)
			req.SetPathValue("id", tt.userID)

			w := httptest.NewRecorder()
			handler.GetUser(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetUser() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var response map[string]interface{}
				json.NewDecoder(w.Body).Decode(&response)

				if response["id"] != tt.userID {
					t.Errorf("GetUser() id = %v, want %v", response["id"], tt.userID)
				}
			}
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	userRepo := newMockUserRepo()
	testUser, _ := domain.NewUser("test@example.com", "Test User")
	testUser.ID = "test_id"
	userRepo.users["test_id"] = testUser

	userService := application.NewUserService(userRepo)
	weatherService := application.NewWeatherService(newMockWeatherService(), userRepo)
	handler := httpHandler.NewHandler(userService, weatherService)

	tests := []struct {
		name       string
		userID     string
		body       map[string]string
		wantStatus int
	}{
		{
			name:   "update name",
			userID: "test_id",
			body: map[string]string{
				"name": "Updated Name",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "update email",
			userID: "test_id",
			body: map[string]string{
				"email": "newemail@example.com",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "non-existing user",
			userID: "non_existing",
			body: map[string]string{
				"name": "New Name",
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("PUT", "/api/users/"+tt.userID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", tt.userID)

			w := httptest.NewRecorder()
			handler.UpdateUser(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("UpdateUser() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	userRepo := newMockUserRepo()
	testUser, _ := domain.NewUser("test@example.com", "Test User")
	testUser.ID = "test_id"
	userRepo.users["test_id"] = testUser

	userService := application.NewUserService(userRepo)
	weatherService := application.NewWeatherService(newMockWeatherService(), userRepo)
	handler := httpHandler.NewHandler(userService, weatherService)

	tests := []struct {
		name       string
		userID     string
		wantStatus int
	}{
		{
			name:       "existing user",
			userID:     "test_id",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "non-existing user",
			userID:     "non_existing",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/users/"+tt.userID, nil)
			req.SetPathValue("id", tt.userID)

			w := httptest.NewRecorder()
			handler.DeleteUser(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("DeleteUser() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_ListUsers(t *testing.T) {
	userRepo := newMockUserRepo()
	for i := 0; i < 3; i++ {
		user, _ := domain.NewUser(string(rune('a'+i))+"@example.com", "User "+string(rune('A'+i)))
		user.ID = string(rune('1' + i))
		userRepo.users[user.ID] = user
	}

	userService := application.NewUserService(userRepo)
	weatherService := application.NewWeatherService(newMockWeatherService(), userRepo)
	handler := httpHandler.NewHandler(userService, weatherService)

	req := httptest.NewRequest("GET", "/api/users?limit=10&offset=0", nil)
	w := httptest.NewRecorder()

	handler.ListUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("ListUsers() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response []map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if len(response) != 3 {
		t.Errorf("ListUsers() returned %d users, want 3", len(response))
	}
}

func TestHandler_GetWeather(t *testing.T) {
	userRepo := newMockUserRepo()
	userService := application.NewUserService(userRepo)
	weatherService := application.NewWeatherService(newMockWeatherService(), userRepo)
	handler := httpHandler.NewHandler(userService, weatherService)

	tests := []struct {
		name       string
		city       string
		wantStatus int
	}{
		{
			name:       "existing city",
			city:       "London",
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing city",
			city:       "NonExisting",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "missing city parameter",
			city:       "",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/weather"
			if tt.city != "" {
				url += "?city=" + tt.city
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler.GetWeather(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetWeather() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var response map[string]interface{}
				json.NewDecoder(w.Body).Decode(&response)

				if response["city"] != tt.city {
					t.Errorf("GetWeather() city = %v, want %v", response["city"], tt.city)
				}
			}
		})
	}
}

func TestHandler_Health(t *testing.T) {
	userRepo := newMockUserRepo()
	userService := application.NewUserService(userRepo)
	weatherService := application.NewWeatherService(newMockWeatherService(), userRepo)
	handler := httpHandler.NewHandler(userService, weatherService)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["status"] != "healthy" {
		t.Errorf("Health() status = %v, want healthy", response["status"])
	}
}
