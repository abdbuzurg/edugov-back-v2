package handlers

import (
	"edugov-back-v2/internal/dto"
	"net/http"
)

// Register handles user registration requests.
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.Register(r.Context(), &req); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login handles user login requests.
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.Login(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

// RefreshToken rotates auth tokens using a refresh token.
func (h *Handlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

// Logout removes active sessions for the supplied refresh token.
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.Logout(r.Context(), &req); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
