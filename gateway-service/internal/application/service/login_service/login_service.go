package login_service

import (
	"context"
	"errors"
	"fmt"
	"gateway-service/internal/application/dto"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FetchUserByEmail(ctx context.Context, userEmail string) (dto.User, error)
}

func (s *service) Login(ctx context.Context, loginingUser *dto.User) (user dto.User, err error) {
	user, err = s.userRepository.FetchUserByEmail(ctx, loginingUser.Email)
	if err != nil {
		return user, errors.New("cannot create a user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginingUser.Password))
	if err != nil {
		return user, errors.New("login error")
	}

	// generate JWT token
	expirationTime := time.Now().Add(30 * time.Minute).Unix()
	jwtClaims := jwt.MapClaims{
		"user id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	jwtToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		fmt.Println(err)
		return user, errors.New("login error")
	}

	user.JWTToken = jwtToken
	return user, nil
}

func New(userRepository UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}

type service struct {
	userRepository UserRepository
}
