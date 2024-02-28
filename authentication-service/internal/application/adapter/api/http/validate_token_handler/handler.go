package validate_token_handler

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

type ValidateTokenService interface {
	ValidateToken(ctx context.Context, validateTokenRequest *dto.ValidateTokenRequest) error
}

func (h *handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	var validateTokenRequest dto.ValidateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&validateTokenRequest); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("the decoding of a valid login request failed", zap.Error(err))
		return
	}

	logger = logger.With(
		zap.String("user_id", validateTokenRequest.UserID.String()),
		zap.String("session_id", validateTokenRequest.UserID.String()),
		zap.Bool("is_bloked", validateTokenRequest.IsBlocked),
		zap.String("access_token", validateTokenRequest.AccessToken),
		zap.Time("access_token_expires_at", validateTokenRequest.AccessTokenExpiresAt),
		zap.String("refresh_token", validateTokenRequest.RefreshToken),
		zap.Time("refresh_token_expires_at", validateTokenRequest.RefreshTokenExpiresAt),
	)

	err := h.validateTokenService.ValidateToken(ctx, &validateTokenRequest)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("validate token is failed", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func New(validateTokenService ValidateTokenService) *handler {
	return &handler{
		validateTokenService: validateTokenService,
	}
}

type handler struct {
	validateTokenService ValidateTokenService
}
