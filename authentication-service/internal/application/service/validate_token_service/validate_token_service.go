package validate_token_service

import (
	"authentication-service/internal/application/dto"
	"authentication-service/internal/application/helper/logging"
	"errors"
	"fmt"
	"os"

	"context"

	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	SessionRepository interface {
		FetchSessionByUserID(ctx context.Context, userID uuid.UUID) (dto.Session, error)
	}
)

func (s *service) ValidateToken(ctx context.Context, accessToken *dto.ValidateToken) error {
	logger := logging.LoggerFromContext(ctx)

	session, err := s.sessionRepository.FetchSessionByUserID(ctx, accessToken.UserID)
	if err != nil {
		logger.Error("fetching the session in the database failed", zap.Error(err))
		return errors.New("cannot fetch the session")
	}

	token, err := jwt.Parse(accessToken.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			logger.Error("token method is incorrect", zap.Error(err))
			return nil, err
		}

		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		logger.Error("token parsing is failing", zap.Error(err))
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			logger.Error("token is expired", zap.Error(err))
			return err
		}

		if session.UserID != claims["user_id"] {
			logger.Error("user from the token doesn't match", zap.Error(err))
			return err
		}
	} else {
		err := errors.New("token claims error")
		logger.Error("token claims error", zap.Error(err))
		return err
	}

	return nil
}

func New(sessionRepository SessionRepository) *service {
	return &service{
		sessionRepository: sessionRepository,
	}
}

type service struct {
	sessionRepository SessionRepository
}
