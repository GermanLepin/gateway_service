package create_user_service

import (
	"context"

	"errors"

	"gateway-service/internal/application/dto"
)

type UserRepository interface {
	CreateUserById(ctx context.Context, user dto.CretaeUserRequest) error
}

func (s *service) CreateUser(ctx context.Context, user dto.CretaeUserRequest) error {
	if err := s.userRepository.CreateUserById(ctx, user); err != nil {
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
