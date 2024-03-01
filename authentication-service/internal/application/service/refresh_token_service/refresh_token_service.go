package refresh_token_service

import (
	"authentication-service/internal/application/dto"
	"authentication-service/internal/application/helper/logging"
	"errors"
	"fmt"
	"os"
	"time"

	"context"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	SessionRepository interface {
		FetchSessionByUserID(ctx context.Context, userID uuid.UUID) (dto.Session, error)
	}

	CreateSessionService interface {
		CreateSession(ctx context.Context, createSession dto.CreateSession) (session dto.Session, err error)
	}
)

func (s *service) RefreshToken(
	ctx context.Context,
	refreshToken *dto.RefreshToken,
) (session dto.Session, err error) {
	logger := logging.LoggerFromContext(ctx)

	currentSession, err := s.sessionRepository.FetchSessionByUserID(ctx, refreshToken.UserID)
	if err != nil {
		logger.Error("fetching the session in the database failed", zap.Error(err))
		return session, errors.New("cannot fetch the session")
	}

	if currentSession.RefreshToken != refreshToken.RefreshToken {
		logger.Error("refresh tokens do not match", zap.Error(err))
		return session, errors.New("refresh tokens do not match")
	}

	token, err := jwt.Parse(refreshToken.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			logger.Error("token method is incorrect", zap.Error(err))
			return nil, err
		}

		return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
	})
	if err != nil {
		logger.Error("token parsing is failing", zap.Error(err))
		return session, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			logger.Error("token is expired", zap.Error(err))
			return session, err
		}
	} else {
		err := errors.New("token claims error")
		logger.Error("token claims error", zap.Error(err))
		return session, err
	}

	createSession := dto.CreateSession{
		UserID:    refreshToken.UserID,
		UserIP:    refreshToken.UserIP,
		SessionID: refreshToken.SessionID,
	}

	// if everything is on track, we should create a new session and return it to the gateway-service
	session, err = s.createSessionService.CreateSession(ctx, createSession)
	if err != nil {
		logger.Error("creating a new session in the database failed", zap.Error(err))
		return session, errors.New("cannot cretae a new session")
	}

	return session, nil
}

func New(
	sessionRepository SessionRepository,
	createSessionService CreateSessionService,
) *service {
	return &service{
		sessionRepository:    sessionRepository,
		createSessionService: createSessionService,
	}
}

type service struct {
	sessionRepository    SessionRepository
	createSessionService CreateSessionService
}
