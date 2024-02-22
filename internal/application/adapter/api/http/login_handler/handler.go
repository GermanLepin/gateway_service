package login_handler

import (
	"context"
	"encoding/json"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"

	"net/http"
	"time"

	"go.uber.org/zap"
)

type LoginService interface {
	Login(ctx context.Context, loginingUser *dto.User) (dto.User, error)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("decoding of login request is failed", zap.Error(err))
		return
	}

	logger = logger.With(zap.String("email", loginRequest.UserID.String()))
	loginingUser := &dto.User{
		ID:       loginRequest.UserID,
		Password: loginRequest.Password,
	}

	user, err := h.loginService.Login(ctx, loginingUser)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("login is failed", zap.Error(err))
		return
	}

	cookieWithAccessToken := http.Cookie{
		Name:     "Access_token",
		Value:    user.AccessToken,
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
		Value:    user.AccessToken,
		MaxAge:   3600 * 24 * 30,
		Path:     "",
		Domain:   "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookieWithRefreshToken)

	descriptionMessage := "user logined successfully"
	loginResponse := dto.LoginResponse{
		UserID:       user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Phone:        user.Phone,
		UserType:     user.UserType,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Message:      descriptionMessage,
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&loginResponse)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("encoding of create user responce is failed", zap.Error(err))
		return
	}
}

func New(loginService LoginService) *handler {
	return &handler{
		loginService: loginService,
	}
}

type handler struct {
	loginService LoginService
}
