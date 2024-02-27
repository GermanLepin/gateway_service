package dto

import "github.com/google/uuid"

type LoginRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}
