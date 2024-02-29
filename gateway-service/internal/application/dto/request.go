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

type LoginRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type ValidateTokenRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type RefreshTokenRequest struct {
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
