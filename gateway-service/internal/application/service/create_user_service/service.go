package create_user_service

import (
	"context"

	"errors"

	"gateway-service/internal/application/dto"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *dto.User) error
}

func (s *service) CreateUser(ctx context.Context, user *dto.User) error {
	if err := s.userRepository.CreateUser(ctx, user); err != nil {
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
