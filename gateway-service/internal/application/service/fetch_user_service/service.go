package fetch_user_service

import (
	"context"
	"fmt"

	"gateway-service/internal/application/dto"

	"github.com/google/uuid"
)

type UserRepository interface {
	FetchUserById(ctx context.Context, userID uuid.UUID) (dto.User, error)
}

func (s *service) FetchBalanceInfo(ctx context.Context, userID uuid.UUID) (dto.User, error) {
	user, err := s.userRepository.FetchUserById(ctx, userID)
	if err != nil {
		return user, fmt.Errorf("cannot fetch the user: %s", userID)
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
