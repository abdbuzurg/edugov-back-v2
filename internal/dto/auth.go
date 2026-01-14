// Package dto defines transport shapes for HTTP requests and responses.
package dto

// ---- REQUEST DTOs ----
// RegisterRequest represents user registration input.
type RegisterRequest struct {
	Tin      string `json:"tin" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthRequest represents login credentials.
type AuthRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LogoutRequest represents a logout request with a refresh token.
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// RefreshTokenRequest represents a token refresh request.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// ---- RESPONSE DTOs ----
// AuthResponse represents the auth token response payload.
type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UniqueID     string `json:"uniqueID"`
}

// MeResponse exposes basic identity info.
type MeResponse struct {
	UniqueID string `json:"uniqueID"`
}
