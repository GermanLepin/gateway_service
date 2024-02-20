package login_service

import (
	"context"
	"errors"
	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/logging"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FetchUserByEmail(ctx context.Context, email string) (dto.User, error)
}

func (s *service) Login(ctx context.Context, loginingUser *dto.User) (dto.User, error) {
	logger := logging.LoggerFromContext(ctx)

	user, err := s.userRepository.FetchUserByEmail(ctx, loginingUser.Email)
	if err != nil {
		logger.Error(
			"fetching user by email in database is failed",
			zap.Error(err),
		)
		return user, errors.New("cannot login the user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginingUser.Password))
	if err != nil {
		logger.Error(
			"received password is incorect",
			zap.Error(err),
		)
		return user, errors.New("login error")
	}

	jwtToken, err := s.jwtTokenGenerator(ctx, &user)
	if err != nil {
		logger.Error(
			"jwt token generation is failed",
			zap.Error(err),
		)
		return user, errors.New("login error")
	}

	user.JWTToken = jwtToken
	return user, nil
}

func (s *service) jwtTokenGenerator(ctx context.Context, user *dto.User) (jwtToken string, err error) {
	logger := logging.LoggerFromContext(ctx)

	// generate JWT token
	expirationTime := time.Now().Add(30 * time.Minute).Unix()
	jwtClaims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"first_name": user.FirstName,
		"email":      user.Email,
		"exp":        expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	jwtToken, err = token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		logger.Error(
			"jwt token generation is failed",
			zap.Error(err),
		)
		return "", errors.New("login error")
	}

	return jwtToken, nil
}

func New(userRepository UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}

type service struct {
	userRepository UserRepository
}
