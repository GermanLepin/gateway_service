package dto

import (
	"time"

	"github.com/google/uuid"
)

type LoginUserRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	UserIP   string    `json:"user_ip"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type ValidateTokenRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	UserIP      string    `json:"user_ip"`
	SessionID   uuid.UUID `json:"session_id"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type RefreshTokenRequest struct {
	UserID       uuid.UUID `json:"user_id"`
	UserIP       string    `json:"user_ip"`
	SessionID    uuid.UUID `json:"session_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
