package delete_user_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserRepository interface {
	DeleteUserById(ctx context.Context, userID uuid.UUID) error
}

func (s *service) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepository.DeleteUserById(ctx, userID); err != nil {
		return fmt.Errorf("cannot delete the user: %s", userID)
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
