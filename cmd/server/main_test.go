package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apiClient "github.com/leinonen/hexagonal-architecture-go/internal/adapters/api"
	httpHandler "github.com/leinonen/hexagonal-architecture-go/internal/adapters/http"
	"github.com/leinonen/hexagonal-architecture-go/internal/adapters/repository/memory"
	"github.com/leinonen/hexagonal-architecture-go/internal/application"
)

func setupTestServer() *httptest.Server {
	userRepo := memory.NewUserRepository()
	weatherClient := apiClient.NewWeatherClient("test-key")

	userService := application.NewUserService(userRepo)
	weatherService := application.NewWeatherService(weatherClient, userRepo)

	handler := httpHandler.NewHandler(userService, weatherService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return httptest.NewServer(mux)
}

func TestIntegration_UserCRUD(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	t.Run("Create User", func(t *testing.T) {
		reqBody := map[string]string{
			"email": "integration@test.com",
			"name":  "Integration Test",
		}
		body, _ := json.Marshal(reqBody)

		resp, err := client.Post(server.URL+"/api/users", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Create user status = %v, want %v", resp.StatusCode, http.StatusCreated)
		}

		var createResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&createResp)

		if createResp["id"] == "" {
			t.Errorf("Create user response missing ID")
		}

		userID := createResp["id"].(string)

		t.Run("Get Created User", func(t *testing.T) {
			resp, err := client.Get(server.URL + "/api/users/" + userID)
			if err != nil {
				t.Fatalf("Failed to get user: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Get user status = %v, want %v", resp.StatusCode, http.StatusOK)
			}

			var getResp map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&getResp)

			if getResp["email"] != "integration@test.com" {
				t.Errorf("Get user email = %v, want integration@test.com", getResp["email"])
			}
		})

		t.Run("Update User", func(t *testing.T) {
			updateBody := map[string]string{
				"name": "Updated Integration Test",
			}
			body, _ := json.Marshal(updateBody)

			req, _ := http.NewRequest("PUT", server.URL+"/api/users/"+userID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to update user: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Update user status = %v, want %v", resp.StatusCode, http.StatusOK)
			}

			var updateResp map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&updateResp)

			if updateResp["name"] != "Updated Integration Test" {
				t.Errorf("Update user name = %v, want Updated Integration Test", updateResp["name"])
			}
		})

		t.Run("List Users", func(t *testing.T) {
			resp, err := client.Get(server.URL + "/api/users")
			if err != nil {
				t.Fatalf("Failed to list users: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("List users status = %v, want %v", resp.StatusCode, http.StatusOK)
			}

			var listResp []map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&listResp)

			if len(listResp) == 0 {
				t.Errorf("List users returned empty list")
			}
		})

		t.Run("Delete User", func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", server.URL+"/api/users/"+userID, nil)

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to delete user: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("Delete user status = %v, want %v", resp.StatusCode, http.StatusNoContent)
			}

			resp, _ = client.Get(server.URL + "/api/users/" + userID)
			if resp.StatusCode != http.StatusNotFound {
				t.Errorf("Get deleted user status = %v, want %v", resp.StatusCode, http.StatusNotFound)
			}
		})
	})
}

func TestIntegration_ErrorHandling(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	tests := []struct {
		name       string
		method     string
		path       string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "Get non-existing user",
			method:     "GET",
			path:       "/api/users/non_existing",
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "Create user without email",
			method: "POST",
			path:   "/api/users",
			body: map[string]string{
				"name": "No Email",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "Create user without name",
			method: "POST",
			path:   "/api/users",
			body: map[string]string{
				"email": "no-name@test.com",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Update non-existing user",
			method:     "PUT",
			path:       "/api/users/non_existing",
			body:       map[string]string{"name": "New Name"},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Delete non-existing user",
			method:     "DELETE",
			path:       "/api/users/non_existing",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				body, _ := json.Marshal(tt.body)
				req, _ = http.NewRequest(tt.method, server.URL+tt.path, bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(tt.method, server.URL+tt.path, nil)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("%s status = %v, want %v", tt.name, resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestIntegration_DuplicateEmail(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	user1 := map[string]string{
		"email": "duplicate@test.com",
		"name":  "First User",
	}
	body1, _ := json.Marshal(user1)

	resp1, err := client.Post(server.URL+"/api/users", "application/json", bytes.NewReader(body1))
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}
	resp1.Body.Close()

	if resp1.StatusCode != http.StatusCreated {
		t.Fatalf("First user creation failed: %v", resp1.StatusCode)
	}

	user2 := map[string]string{
		"email": "duplicate@test.com",
		"name":  "Second User",
	}
	body2, _ := json.Marshal(user2)

	resp2, err := client.Post(server.URL+"/api/users", "application/json", bytes.NewReader(body2))
	if err != nil {
		t.Fatalf("Failed to create second user: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusConflict {
		t.Errorf("Duplicate email status = %v, want %v", resp2.StatusCode, http.StatusConflict)
	}
}

func TestIntegration_HealthCheck(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Health check status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	var health map[string]string
	json.NewDecoder(resp.Body).Decode(&health)

	if health["status"] != "healthy" {
		t.Errorf("Health check status = %v, want healthy", health["status"])
	}
}
