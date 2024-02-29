package refresh_token_handler

import (
	"context"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"

	"net/http"
	"time"

	"go.uber.org/zap"
)

type RefreshTokenService interface {
	RefreshTokenRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (session dto.Session, err error)
}

func (h *handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	session, err := h.refreshTokenService.RefreshTokenRequest(ctx, w, r)
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

	loginResponse := dto.RefreshTokenResponse{
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

func New(refreshTokenService RefreshTokenService) *handler {
	return &handler{
		refreshTokenService: refreshTokenService,
	}
}

type handler struct {
	refreshTokenService RefreshTokenService
}
