package repository

import (
	"auth-service/internal/application/dto"

	"context"
	"database/sql"
)

func (s *sessionRepository) SaveSession(ctx context.Context, session *dto.Session) error {
	queryString := `
		insert into service.sessions(id, user_id, is_blocked, refresh_token, expires_at) 
		values ($1,$2,$3,$4,$5)
	;`

	err := s.db.QueryRow(queryString,
		session.ID,
		session.UserID,
		session.IsBlocked,
		session.RefreshToken,
		session.RefreshTokenExpiresAt,
	)
	if err != nil {
		return err.Err()
	}

	return nil
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *sessionRepository {
	return &sessionRepository{
		db: db,
	}
}
