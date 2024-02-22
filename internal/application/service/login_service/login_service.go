package login_service

import (
	"context"
	"errors"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/logging"

	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FetchUserById(ctx context.Context, userId uuid.UUID) (dto.User, error)
}

func (s *service) Login(ctx context.Context, loginingUser *dto.User) (dto.User, error) {
	logger := logging.LoggerFromContext(ctx)

	user, err := s.userRepository.FetchUserById(ctx, loginingUser.ID)
	if err != nil {
		logger.Error("fetching user by email in database is failed", zap.Error(err))
		return user, errors.New("cannot login the user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginingUser.Password))
	if err != nil {
		logger.Error("received password is incorect", zap.Error(err))
		return user, errors.New("login error")
	}

	accessToken, err := s.generateAccessToken(ctx, &user)
	if err != nil {
		logger.Error("access jwt token generation is failed", zap.Error(err))
		return user, errors.New("login error")
	}

	refreshToken, err := s.generateRefreshToken(ctx, &user)
	if err != nil {
		logger.Error("refresh jwt token generation is failed", zap.Error(err))
		return user, errors.New("login error")
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	return user, nil
}

func (s *service) generateAccessToken(ctx context.Context, user *dto.User) (jwtToken string, err error) {
	logger := logging.LoggerFromContext(ctx)

	// generate JWT token
	expirationTime := time.Now().Add(15 * time.Minute).Unix()
	jwtClaims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"first_name": user.FirstName,
		"email":      user.Email,
		"exp":        expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	jwtToken, err = token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		logger.Error("access jwt token generation is failed", zap.Error(err))
		return "", errors.New("login error")
	}

	return jwtToken, nil
}

func (s *service) generateRefreshToken(ctx context.Context, user *dto.User) (jwtToken string, err error) {
	logger := logging.LoggerFromContext(ctx)

	// generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	jwtClaims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"first_name": user.FirstName,
		"email":      user.Email,
		"exp":        expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	jwtToken, err = token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		logger.Error("refresh jwt token generation is failed", zap.Error(err))
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
