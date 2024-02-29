package repository

import (
	"authentication-service/internal/application/dto"

	"context"
	"database/sql"

	"github.com/google/uuid"
)

func (s *sessionRepository) SaveSession(ctx context.Context, session *dto.Session) error {
	queryString := `
		insert into service.sessions(id, user_id, is_blocked, refresh_token, refresh_token_expires_at)
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

func (s *sessionRepository) FetchSessionByUserID(ctx context.Context, userID uuid.UUID) (session dto.Session, err error) {
	queryString := `
		select * from service.sessions
		where user_id = $1
		order by created_at desc;
	;`

	err = s.db.QueryRow(queryString, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.IsBlocked,
		&session.RefreshToken,
		&session.RefreshTokenExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return session, err
	}

	return session, nil
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *sessionRepository {
	return &sessionRepository{
		db: db,
	}
}
