package create_session_service

import (
	"auth-service/internal/application/dto"
	"auth-service/internal/application/helper/logging"

	"context"
	"errors"

	"go.uber.org/zap"
)

type SessionRepository interface {
	SaveSession(ctx context.Context, session *dto.Session) error
}

func (s *service) CreateSession(ctx context.Context, session *dto.Session) error {
	logger := logging.LoggerFromContext(ctx)

	if err := s.sessionRepository.SaveSession(ctx, session); err != nil {
		logger.Error("creation session in database is failed", zap.Error(err))
		return errors.New("cannot create a session")
	}

	return nil
}

func New(sessionRepository SessionRepository) *service {
	return &service{
		sessionRepository: sessionRepository,
	}
}

type service struct {
	sessionRepository SessionRepository
}
