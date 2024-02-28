package dto

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserID                uuid.UUID
	ID                    uuid.UUID
	IsBlocked             bool
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}
