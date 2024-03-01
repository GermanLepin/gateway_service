package login_user_handler

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

type CreateSessionService interface {
	CreateSession(ctx context.Context, createSession dto.CreateSession) (session dto.Session, err error)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	var loginUserRequest dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginUserRequest); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("the decoding of the logging request failed", zap.Error(err))
		return
	}

	logger = logger.With(
		zap.String("user_id", loginUserRequest.UserID.String()),
		zap.String("user_ip", loginUserRequest.UserIP),
	)

	createSession := dto.CreateSession{
		UserID: loginUserRequest.UserID,
		UserIP: loginUserRequest.UserIP,
	}

	session, err := h.createSessionService.CreateSession(ctx, createSession)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("login is failed", zap.Error(err))
		return
	}

	loginResponse := dto.LoginResponse{
		SessionID:             session.ID,
		IsBlocked:             session.IsBlocked,
		AccessToken:           session.AccessToken,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt,
		RefreshToken:          session.RefreshToken,
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
		UserID:                session.UserID,
	}

	if err = jsonwrapper.WriteJSON(w, http.StatusAccepted, loginResponse); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("cannot send a login response", zap.Error(err))
		return
	}
}

func New(createSessionService CreateSessionService) *handler {
	return &handler{
		createSessionService: createSessionService,
	}
}

type handler struct {
	createSessionService CreateSessionService
}
