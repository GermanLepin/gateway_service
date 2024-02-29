package fetch_user_service

import (
	"context"
	"fmt"
	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/logging"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	UserRepository interface {
		FetchUserById(ctx context.Context, userId uuid.UUID) (dto.User, error)
	}

	ValidateTokenService interface {
		ValidateTokenRequest(ctx context.Context, r *http.Request, userID uuid.UUID) error
	}
)

func (s *service) FetchUser(
	ctx context.Context,
	r *http.Request,
	userId uuid.UUID,
) (user dto.User, err error) {
	logger := logging.LoggerFromContext(ctx)

	if err = s.validateTokenService.ValidateTokenRequest(ctx, r, userId); err != nil {
		logger.Error("validate token request is failed", zap.Error(err))
		return user, fmt.Errorf("validate token request is failed: %s", userId)
	}

	user, err = s.userRepository.FetchUserById(ctx, userId)
	if err != nil {
		logger.Error("fetching user by email in database is failed", zap.Error(err))
		return user, fmt.Errorf("cannot fetch the user: %s", userId)
	}

	return user, nil
}

func New(
	userRepository UserRepository,
	validateTokenService ValidateTokenService,
) *service {
	return &service{
		userRepository:       userRepository,
		validateTokenService: validateTokenService,
	}
}

type service struct {
	userRepository       UserRepository
	validateTokenService ValidateTokenService
}
