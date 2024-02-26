package create_user_service

import (
	"context"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/logging"

	"errors"

	"go.uber.org/zap"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user *dto.User) error
}

func (s *service) CreateUser(ctx context.Context, user *dto.User) error {
	logger := logging.LoggerFromContext(ctx)

	if err := s.userRepository.SaveUser(ctx, user); err != nil {
		logger.Error("creation user in database is failed", zap.Error(err))
		return errors.New("cannot create a user")
	}

	return nil
}

func New(userRepository UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}

type service struct {
	userRepository UserRepository
}
