package dto

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                    uuid.UUID
	UserID                uuid.UUID
	IsBlocked             bool
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}
