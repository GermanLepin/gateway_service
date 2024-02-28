package refresh_token_service

import (
	"authentication-service/internal/application/dto"

	"context"
)

func (s *service) RefreshToken(
	ctx context.Context,
	refreshTokenRequest *dto.RefreshTokenRequest,
) (session dto.Session, err error) {

	// TODO finish refreshToken
	return session, nil
}

func New() *service {
	return &service{}
}

type service struct{}
