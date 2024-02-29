package dto

import (
	"time"

	"github.com/google/uuid"
)

type ValidateToken struct {
	UserID      uuid.UUID
	AccessToken string
	ExpiresAt   time.Time
}

type RefreshToken struct {
	UserID       uuid.UUID
	RefreshToken string
	ExpiresAt    time.Time
}
