package login_service

import (
	"authentication-service/internal/application/dto"
	"authentication-service/internal/application/helper/logging"

	"os"

	"context"
	"errors"

	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateSessionService interface {
	CreateSession(ctx context.Context, session *dto.Session) error
}

func (s *service) Login(
	ctx context.Context,
	loginRequest *dto.LoginRequest,
) (session dto.Session, err error) {
	logger := logging.LoggerFromContext(ctx)

	accessTokenExpirationTime := time.Now().Add(15 * time.Minute)
	accessToken, err := s.generateAccessToken(ctx, loginRequest, accessTokenExpirationTime)
	if err != nil {
		logger.Error("access jwt token generation is failed", zap.Error(err))
		return session, errors.New("login error")
	}

	refreshTokenExpirationTime := time.Now().Add(72 * time.Hour)
	refreshToken, err := s.generateRefreshToken(ctx, loginRequest, refreshTokenExpirationTime)
	if err != nil {
		logger.Error("refresh jwt token generation is failed", zap.Error(err))
		return session, errors.New("login error")
	}

	session = dto.Session{
		ID:                    uuid.New(),
		IsBlocked:             false,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpirationTime,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpirationTime,
		UserID:                loginRequest.UserID,
	}

	err = s.createSessionService.CreateSession(ctx, &session)
	if err != nil {
		logger.Error("fetching user by email in database is failed", zap.Error(err))
		return session, errors.New("cannot login the user")
	}

	return session, nil
}

func (s *service) generateAccessToken(
	ctx context.Context,
	loginRequest *dto.LoginRequest,
	accessTokenExpirationTime time.Time,
) (jwtToken string, err error) {
	logger := logging.LoggerFromContext(ctx)

	// generate access JWT token
	jwtClaims := jwt.MapClaims{
		"user_id": loginRequest.UserID,
		"email":   loginRequest.Email,
		"exp":     accessTokenExpirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	jwtToken, err = token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		logger.Error("access jwt token generation is failed", zap.Error(err))
		return "", errors.New("login error")
	}

	return jwtToken, nil
}

func (s *service) generateRefreshToken(
	ctx context.Context,
	loginRequest *dto.LoginRequest,
	refreshTokenExpirationTime time.Time,
) (jwtToken string, err error) {
	logger := logging.LoggerFromContext(ctx)

	// generate refresh JWT token
	jwtClaims := jwt.MapClaims{
		"user_id": loginRequest.UserID,
		"email":   loginRequest.Email,
		"exp":     refreshTokenExpirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	jwtToken, err = token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		logger.Error("refresh jwt token generation is failed", zap.Error(err))
		return "", errors.New("login error")
	}

	return jwtToken, nil
}

func New(createSessionService CreateSessionService) *service {
	return &service{
		createSessionService: createSessionService,
	}
}

type service struct {
	createSessionService CreateSessionService
}
