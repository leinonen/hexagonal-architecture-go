package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leinonen/hexagonal-architecture-go/internal/application"
	"github.com/leinonen/hexagonal-architecture-go/internal/errors"
)

type Handler struct {
	userService    *application.UserService
	weatherService *application.WeatherService
}

func NewHandler(userService *application.UserService, weatherService *application.WeatherService) *Handler {
	return &Handler{
		userService:    userService,
		weatherService: weatherService,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("GET /api/users", h.ListUsers)
	mux.HandleFunc("GET /api/users/{id}", h.GetUser)
	mux.HandleFunc("PUT /api/users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", h.DeleteUser)

	mux.HandleFunc("GET /api/weather", h.GetWeather)
	mux.HandleFunc("GET /api/users/{id}/weather", h.GetUserWeather)

	mux.HandleFunc("GET /health", h.Health)
}

func (h *Handler) respondWithError(w http.ResponseWriter, err error) {
	var status int
	var message string

	switch {
	case errors.IsNotFound(err):
		status = http.StatusNotFound
		message = err.Error()
	case errors.IsValidation(err):
		status = http.StatusBadRequest
		message = err.Error()
	case errors.IsConflict(err):
		status = http.StatusConflict
		message = err.Error()
	case errors.IsUnauthorized(err):
		status = http.StatusUnauthorized
		message = err.Error()
	case errors.IsExternalService(err):
		status = http.StatusServiceUnavailable
		message = "External service unavailable"
	default:
		status = http.StatusInternalServerError
		message = "Internal server error"
		log.Printf("Internal error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
