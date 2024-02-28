package dto

import (
	"time"

	"github.com/google/uuid"
)

type LoginResponse struct {
	UserID                uuid.UUID `json:"user_id"`
	SessionID             uuid.UUID `json:"session_id"`
	IsBlocked             bool      `json:"is_bloked"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type RefreshTokenResponce struct {
	UserID                uuid.UUID `json:"user_id"`
	SessionID             uuid.UUID `json:"session_id"`
	IsBlocked             bool      `json:"is_bloked"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
