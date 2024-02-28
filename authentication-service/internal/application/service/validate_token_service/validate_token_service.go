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

func (s *service) ValidateToken(ctx context.Context, validateTokenRequest *dto.ValidateTokenRequest) error {
	logger := logging.LoggerFromContext(ctx)

	token, err := jwt.Parse(validateTokenRequest.AccessToken, func(token *jwt.Token) (interface{}, error) {
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

			// TODO refresh token
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
