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
	"go.uber.org/zap"
)

func (s *service) ValidateToken(ctx context.Context, validateToken *dto.ValidateToken) error {
	logger := logging.LoggerFromContext(ctx)

	// TODO check: do we have a user in db?

	token, err := jwt.Parse(validateToken.AccessToken, func(token *jwt.Token) (interface{}, error) {
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
	} else {
		err := errors.New("token claims error")
		logger.Error("token claims error", zap.Error(err))
		return err
	}

	return nil
}

func New() *service {
	return &service{}
}

type service struct{}
