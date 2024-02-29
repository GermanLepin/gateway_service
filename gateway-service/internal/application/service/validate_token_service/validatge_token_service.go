package validate_token_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) ValidateTokenRequest(ctx context.Context, r *http.Request, userID uuid.UUID) error {
	logger := logging.LoggerFromContext(ctx)

	tokenString := r.Header.Get("access_token")
	expiresAtString := r.Header.Get("access_token_expires_at")
	expiresAt, err := time.Parse(time.RFC3339Nano, expiresAtString)
	if err != nil {
		logger.Error("time parcing is failed", zap.Error(err))
		return err
	}

	validateTokenRequest := &dto.ValidateTokenRequest{
		UserID:      userID,
		AccessToken: tokenString,
		ExpiresAt:   expiresAt,
	}

	jsonData, err := json.MarshalIndent(validateTokenRequest, "", "\t")
	if err != nil {
		logger.Error("validate token request marshalling is failed", zap.Error(err))
		return err
	}

	authenticationServiceURL := "http://authentication-service/validate-token"
	request, err := http.NewRequest(http.MethodPost, authenticationServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("cannot reach out to the authentication-service", zap.Error(err))
		return err
	}
	request.Close = true

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Error("sending an HTTP request is a failure", zap.Error(err))
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := errors.New("invalid access token")
		logger.Error("invalid access token", zap.Error(err))
		return err
	}

	return nil
}

func New() *service {
	return &service{}
}

type service struct {
}
