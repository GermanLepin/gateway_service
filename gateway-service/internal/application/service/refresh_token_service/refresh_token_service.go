package refresh_token_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) RefreshTokenRequest(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) (session dto.Session, err error) {
	logger := logging.LoggerFromContext(ctx)

	userID := r.Header.Get("user_id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.Error("user ID parcing is failed", zap.Error(err))
		return session, err
	}

	tokenString := r.Header.Get("refresh_token")

	expiresAtString := r.Header.Get("refresh_token_expires_at")
	expiresAt, err := time.Parse(time.RFC3339Nano, expiresAtString)
	if err != nil {
		logger.Error("time parcing is failed", zap.Error(err))
		return session, err
	}

	refreshTokenRequest := &dto.RefreshTokenRequest{
		UserID:       userUUID,
		RefreshToken: tokenString,
		ExpiresAt:    expiresAt,
	}

	jsonData, err := json.MarshalIndent(refreshTokenRequest, "", "\t")
	if err != nil {
		logger.Error("refresh token request marshalling is failed", zap.Error(err))
		return session, err
	}

	authenticationServiceURL := "http://authentication-service/refresh-token"
	request, err := http.NewRequest(http.MethodPost, authenticationServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("cannot reach out to the authentication-service", zap.Error(err))
		return session, err
	}
	request.Close = true

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Error("sending an HTTP request is a failure", zap.Error(err))
		return session, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := errors.New("invalid refresh token")
		logger.Error("invalid refresh token", zap.Error(err))
		return session, err
	}

	var refreshTokenResponse dto.RefreshTokenResponse
	if err = json.NewDecoder(response.Body).Decode(&refreshTokenResponse); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("the decoding of the login request is failed", zap.Error(err))
		return
	}

	session = dto.Session{
		ID:                    refreshTokenResponse.SessionID,
		IsBlocked:             refreshTokenResponse.IsBlocked,
		AccessToken:           refreshTokenResponse.AccessToken,
		AccessTokenExpiresAt:  refreshTokenResponse.AccessTokenExpiresAt,
		RefreshToken:          refreshTokenResponse.RefreshToken,
		RefreshTokenExpiresAt: refreshTokenResponse.RefreshTokenExpiresAt,
		UserID:                refreshTokenResponse.UserID,
	}

	return session, nil
}

func New() *service {
	return &service{}
}

type service struct {
}
