package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     int    `json:"phone"`
	UserType  string `json:"user_type"`
}

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
