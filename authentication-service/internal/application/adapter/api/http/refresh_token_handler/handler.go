package refresh_token_handler

import (
	"authentication-service/internal/application/dto"
	"authentication-service/internal/application/helper/jsonwrapper"
	"authentication-service/internal/application/helper/logging"

	"context"
	"encoding/json"

	"net/http"
	"time"

	"go.uber.org/zap"
)

type RefreshTokenService interface {
	RefreshToken(ctx context.Context, refreshTokenRequest *dto.RefreshTokenRequest) (dto.Session, error)
}

func (h *handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	var refreshTokenRequest dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshTokenRequest); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("the decoding of the refreshing token request failed", zap.Error(err))
		return
	}

	logger = logger.With(
		zap.String("user_id", refreshTokenRequest.UserID.String()),
		zap.String("session_id", refreshTokenRequest.UserID.String()),
		zap.Bool("is_bloked", refreshTokenRequest.IsBlocked),
		zap.String("access_token", refreshTokenRequest.AccessToken),
		zap.Time("access_token_expires_at", refreshTokenRequest.AccessTokenExpiresAt),
		zap.String("refresh_token", refreshTokenRequest.RefreshToken),
		zap.Time("refresh_token_expires_at", refreshTokenRequest.RefreshTokenExpiresAt),
	)

	session, err := h.refreshTokenService.RefreshToken(ctx, &refreshTokenRequest)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("refresh token is failed", zap.Error(err))
		return
	}

	refreshTokenResponce := dto.RefreshTokenResponce{
		UserID:                session.UserID,
		SessionID:             session.ID,
		IsBlocked:             session.IsBlocked,
		AccessToken:           session.AccessToken,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt,
		RefreshToken:          session.RefreshToken,
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
	}

	if err = jsonwrapper.WriteJSON(w, http.StatusOK, refreshTokenResponce); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("could not send a refresh token response", zap.Error(err))
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
