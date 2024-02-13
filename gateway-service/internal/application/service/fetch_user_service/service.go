package fetch_user_service

import (
	"context"
	"fmt"

	"gateway-service/internal/application/dto"
)

type UserRepository interface {
	FetchUserByEmail(ctx context.Context, userEmail string) (dto.User, error)
}

func (s *service) FetchUser(ctx context.Context, userEmail string) (dto.User, error) {
	user, err := s.userRepository.FetchUserByEmail(ctx, userEmail)
	if err != nil {
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
