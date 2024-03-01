package create_session_service

import (
	"authentication-service/internal/application/dto"
	"authentication-service/internal/application/helper/logging"
	"os"
	"time"

	"context"
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SessionRepository interface {
	SaveSession(ctx context.Context, session *dto.Session) error
}

func (s *service) CreateSession(
	ctx context.Context,
	createSession dto.CreateSession,
) (session dto.Session, err error) {
	logger := logging.LoggerFromContext(ctx)

	sessionID := uuid.New()
	createSession.SessionID = sessionID

	accessTokenExpirationTime := time.Now().Add(15 * time.Minute)
	accessToken, err := s.generateAccessToken(ctx, createSession, accessTokenExpirationTime)
	if err != nil {
		logger.Error("access jwt token generation is failed", zap.Error(err))
		return session, errors.New("login error")
	}

	refreshTokenExpirationTime := time.Now().Add(96 * time.Hour)
	refreshToken, err := s.generateRefreshToken(ctx, createSession, refreshTokenExpirationTime)
	if err != nil {
		logger.Error("refresh jwt token generation is failed", zap.Error(err))
		return session, errors.New("login error")
	}

	session = dto.Session{
		ID:                    uuid.New(),
		UserID:                createSession.UserID,
		IsBlocked:             false,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpirationTime,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpirationTime,
	}

	if err := s.sessionRepository.SaveSession(ctx, &session); err != nil {
		logger.Error("creation session in database is failed", zap.Error(err))
		return session, errors.New("cannot create a new session")
	}

	return session, nil
}

func (s *service) generateAccessToken(
	ctx context.Context,
	createSession dto.CreateSession,
	accessTokenExpirationTime time.Time,
) (jwtToken string, err error) {
	logger := logging.LoggerFromContext(ctx)

	// generate access JWT token
	jwtClaims := jwt.MapClaims{
		"user_id":    createSession.UserID,
		"user_ip":    createSession.UserIP,
		"session_id": createSession.SessionID,
		"exp":        accessTokenExpirationTime.Unix(),
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
	createSession dto.CreateSession,
	refreshTokenExpirationTime time.Time,
) (jwtToken string, err error) {
	logger := logging.LoggerFromContext(ctx)

	// generate access JWT token
	jwtClaims := jwt.MapClaims{
		"user_id":    createSession.UserID,
		"user_ip":    createSession.UserIP,
		"session_id": createSession.SessionID,
		"exp":        refreshTokenExpirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	jwtToken, err = token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		logger.Error("refresh jwt token generation is failed", zap.Error(err))
		return "", errors.New("login error")
	}

	return jwtToken, nil
}

func New(sessionRepository SessionRepository) *service {
	return &service{
		sessionRepository: sessionRepository,
	}
}

type service struct {
	sessionRepository SessionRepository
}
