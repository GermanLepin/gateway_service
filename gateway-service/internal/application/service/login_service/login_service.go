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

type (
	UserRepository interface {
		FetchUserById(ctx context.Context, userId uuid.UUID) (dto.User, error)
	}

	CreateSessionService interface {
		CreateSession(ctx context.Context, session *dto.Session) error
	}
)

func (s *service) Login(
	ctx context.Context,
	loginingUser *dto.User,
) (user dto.User, session dto.Session, err error) {
	logger := logging.LoggerFromContext(ctx)

	user, err = s.userRepository.FetchUserById(ctx, loginingUser.ID)
	if err != nil {
		logger.Error("fetching user by user id in database is failed", zap.Error(err))
		return user, session, errors.New("cannot login the user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginingUser.Password))
	if err != nil {
		logger.Error("received password is incorect", zap.Error(err))
		return user, session, errors.New("login error")
	}

	accessTokenExpirationTime := time.Now().Add(15 * time.Minute)
	accessToken, err := s.generateAccessToken(ctx, &user, accessTokenExpirationTime)
	if err != nil {
		logger.Error("access jwt token generation is failed", zap.Error(err))
		return user, session, errors.New("login error")
	}

	refreshTokenExpirationTime := time.Now().Add(24 * time.Hour)
	refreshToken, err := s.generateRefreshToken(ctx, &user, refreshTokenExpirationTime)
	if err != nil {
		logger.Error("refresh jwt token generation is failed", zap.Error(err))
		return user, session, errors.New("login error")
	}

	session = dto.Session{
		ID:                    uuid.New(),
		UserID:                user.ID,
		IsBlocked:             false,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpirationTime,
	}

	err = s.createSessionService.CreateSession(ctx, &session)
	if err != nil {
		logger.Error("fetching user by email in database is failed", zap.Error(err))
		return user, session, errors.New("cannot login the user")
	}

	user.AccessToken = accessToken
	user.AccessTokenExpiresAt = accessTokenExpirationTime
	user.RefreshToken = refreshToken
	user.RefreshTokenExpiresAt = refreshTokenExpirationTime
	return user, session, nil
}

func (s *service) generateAccessToken(
	ctx context.Context,
	user *dto.User,
	accessTokenExpirationTime time.Time,
) (jwtToken string, err error) {
	logger := logging.LoggerFromContext(ctx)

	// generate access JWT token
	jwtClaims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"first_name": user.FirstName,
		"email":      user.Email,
		"exp":        accessTokenExpirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	jwtToken, err = token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		logger.Error("access jwt token generation is failed", zap.Error(err))
		return "", errors.New("login error")
	}

	return jwtToken, nil
}

func (s *service) generateRefreshToken(
	ctx context.Context,
	user *dto.User,
	refreshTokenExpirationTime time.Time,
) (jwtToken string, err error) {
	logger := logging.LoggerFromContext(ctx)

	// generate refresh JWT token
	jwtClaims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"first_name": user.FirstName,
		"email":      user.Email,
		"exp":        refreshTokenExpirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	jwtToken, err = token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		logger.Error("refresh jwt token generation is failed", zap.Error(err))
		return "", errors.New("login error")
	}

	return jwtToken, nil
}

func New(
	userRepository UserRepository,
	createSessionService CreateSessionService,
) *service {
	return &service{
		userRepository:       userRepository,
		createSessionService: createSessionService,
	}
}

type service struct {
	userRepository       UserRepository
	createSessionService CreateSessionService
}
