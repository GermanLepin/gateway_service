package fetch_user_service

import (
	"context"
	"fmt"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/logging"

	"go.uber.org/zap"
)

type UserRepository interface {
	FetchUserByEmail(ctx context.Context, userEmail string) (dto.User, error)
}

func (s *service) FetchUser(ctx context.Context, userEmail string) (dto.User, error) {
	logger := logging.LoggerFromContext(ctx)

	user, err := s.userRepository.FetchUserByEmail(ctx, userEmail)
	if err != nil {
		logger.Error(
			"fetching user by email in database is failed",
			zap.Error(err),
		)
		return user, fmt.Errorf("cannot fetch the user: %s", userEmail)
	}

	return user, nil
}

func New(userRepository UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}

type service struct {
	userRepository UserRepository
}
