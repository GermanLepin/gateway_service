package login_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FetchUserByEmail(ctx context.Context, email string) (user dto.User, err error)
}

func (s *service) Login(
	ctx context.Context,
	w http.ResponseWriter,
	loginRequest *dto.LoginRequest,
) (session dto.Session, err error) {
	logger := logging.LoggerFromContext(ctx)

	user, err := s.userRepository.FetchUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		logger.Error("fetching user by user id in database is failed", zap.Error(err))
		return session, errors.New("cannot login the user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		logger.Error("received password is incorect", zap.Error(err))
		return session, errors.New("login error")
	}

	loginRequest.UserID = user.ID
	session, err = s.authenticate(ctx, w, *loginRequest)
	if err != nil {
		logger.Error("sending a request to authentication-service is failed", zap.Error(err))
		return session, errors.New("login error")
	}

	cookieWithAccessToken := http.Cookie{
		Name:     "Access_token",
		Value:    session.AccessToken,
		MaxAge:   3600 * 24 * 30,
		Path:     "",
		Domain:   "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookieWithAccessToken)

	cookieWithRefreshToken := http.Cookie{
		Name:     "Refresh_token",
		Value:    session.AccessToken,
		MaxAge:   3600 * 24 * 30,
		Path:     "",
		Domain:   "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookieWithRefreshToken)

	return session, nil
}

func (s *service) authenticate(
	ctx context.Context,
	w http.ResponseWriter,
	loginRequest dto.LoginRequest,
) (session dto.Session, err error) {
	logger := logging.LoggerFromContext(ctx)

	jsonData, err := json.MarshalIndent(loginRequest, "", "\t")
	if err != nil {
		logger.Error("login request marshalling is failed", zap.Error(err))
		return session, err
	}

	authenticationServiceURL := "http://authentication-service/user/login"
	request, err := http.NewRequest(http.MethodPost, authenticationServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("cannot reach out the authentication-service", zap.Error(err))
		return session, err
	}
	request.Close = true

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Error("sending an HTTP request is failed", zap.Error(err))
		return session, err
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		err := errors.New("invalid credentials")
		logger.Error("invalid credentials", zap.Error(err))
		return session, err
	} else if response.StatusCode != http.StatusAccepted {
		err := errors.New("error calling the authentication service")
		logger.Error("error calling the authentication service", zap.Error(err))
		return session, err
	}

	var loginResponse dto.LoginResponse
	if err = json.NewDecoder(response.Body).Decode(&loginResponse); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("decoding of login request is failed", zap.Error(err))
		return
	}

	session = dto.Session{
		ID:                    loginResponse.SessionID,
		IsBlocked:             loginResponse.IsBlocked,
		AccessToken:           loginResponse.AccessToken,
		AccessTokenExpiresAt:  loginResponse.AccessTokenExpiresAt,
		RefreshToken:          loginResponse.RefreshToken,
		RefreshTokenExpiresAt: loginResponse.RefreshTokenExpiresAt,
		UserID:                loginResponse.UserID,
	}

	return session, nil
}

func New(userRepository UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}

type service struct {
	userRepository UserRepository
}
