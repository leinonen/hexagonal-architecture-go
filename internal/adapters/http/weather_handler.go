package http

import (
	"net/http"

	"github.com/leinonen/hexagonal-architecture-go/internal/dto"
	"github.com/leinonen/hexagonal-architecture-go/internal/errors"
)

func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		h.respondWithError(w, errors.NewValidationError("city parameter is required"))
		return
	}

	weather, err := h.weatherService.GetWeather(r.Context(), city)
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	response := dto.ToWeatherResponseDTO(weather)

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) GetUserWeather(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	city := r.URL.Query().Get("city")

	if city == "" {
		h.respondWithError(w, errors.NewValidationError("city parameter is required"))
		return
	}

	weather, err := h.weatherService.GetWeatherForUser(r.Context(), userID, city)
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	response := dto.ToWeatherResponseDTO(weather)

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}
