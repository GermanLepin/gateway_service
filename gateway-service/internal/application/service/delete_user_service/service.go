package delete_user_service

import (
	"context"
	"fmt"

	"gateway-service/gateway-service/internal/application/helper/logging"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserRepository interface {
	DeleteUserById(ctx context.Context, userId uuid.UUID) error
}

func (s *service) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	logger := logging.LoggerFromContext(ctx)

	if err := s.userRepository.DeleteUserById(ctx, userId); err != nil {
		logger.Error("deletion user in database is failed", zap.Error(err))
		return fmt.Errorf("cannot delete the user: %s", userId)
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
