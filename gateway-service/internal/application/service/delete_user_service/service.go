package delete_user_service

import (
	"context"
	"fmt"
	"gateway-service/internal/application/helper/logging"

	"go.uber.org/zap"
)

type UserRepository interface {
	DeleteUserByEmail(ctx context.Context, email string) error
}

func (s *service) DeleteUser(ctx context.Context, email string) error {
	logger := logging.LoggerFromContext(ctx)

	if err := s.userRepository.DeleteUserByEmail(ctx, email); err != nil {
		logger.Error(
			"deletion user in database is failed",
			zap.Error(err),
		)
		return fmt.Errorf("cannot delete the user: %s", email)
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
