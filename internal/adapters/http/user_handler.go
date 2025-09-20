package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/leinonen/hexagonal-architecture-go/internal/dto"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, err)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	response := dto.ToUserResponseDTO(user)

	h.respondWithJSON(w, http.StatusCreated, response)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	response := dto.ToUserResponseDTO(user)

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req dto.UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, err)
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), id, req.Email, req.Name)
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	response := dto.ToUserResponseDTO(user)

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.userService.DeleteUser(r.Context(), id); err != nil {
		h.respondWithError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.userService.ListUsers(r.Context(), limit, offset)
	if err != nil {
		h.respondWithError(w, err)
		return
	}

	response := dto.ToUserResponseDTOs(users)

	h.respondWithJSON(w, http.StatusOK, response)
}
