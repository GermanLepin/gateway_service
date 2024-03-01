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
	Login(ctx context.Context, w http.ResponseWriter, loginUserRequest *dto.LoginUserRequest) (dto.Session, error)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	var loginUserRequest dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginUserRequest); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("the decoding of the login request is failed", zap.Error(err))
		return
	}

	logger = logger.With(zap.String("email", loginUserRequest.Email))

	loginUserRequest.UserIP = r.Header.Get("user_ip")
	session, err := h.loginService.Login(ctx, w, &loginUserRequest)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("login is failed", zap.Error(err))
		return
	}

	cookieWithAccessToken := http.Cookie{
		Name:     "access_token",
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
		Name:     "refresh_token",
		Value:    session.RefreshToken,
		MaxAge:   3600 * 24 * 30,
		Path:     "",
		Domain:   "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookieWithRefreshToken)

	loginResponse := dto.LoginResponse{
		SessionID:             session.ID,
		AccessToken:           session.AccessToken,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt,
		RefreshToken:          session.RefreshToken,
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
		UserID:                session.UserID,
	}

	if err = jsonwrapper.WriteJSON(w, http.StatusOK, loginResponse); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("could not send a login user response", zap.Error(err))
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
