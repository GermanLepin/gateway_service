package dto

import (
	"time"

	"github.com/google/uuid"
)

type ValidateToken struct {
	UserID      uuid.UUID
	UserIP      string
	SessionID   uuid.UUID
	AccessToken string
	ExpiresAt   time.Time
}

type RefreshToken struct {
	UserID       uuid.UUID
	UserIP       string
	SessionID    uuid.UUID
	RefreshToken string
	ExpiresAt    time.Time
}
