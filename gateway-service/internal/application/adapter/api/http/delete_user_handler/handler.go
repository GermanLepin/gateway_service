package delete_user_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type DeleteUserService interface {
	DeleteUser(ctx context.Context, email string) error
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	email := chi.URLParam(r, "email")
	logger = logger.With(zap.String("email", email))

	err := h.deleteUserService.DeleteUser(ctx, email)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"user deletion is failed",
			zap.Error(err),
		)
		return
	}

	descriptionMessage := "user deleted successfully"
	deleteUserResponse := dto.DeleteUserResponse{
		Email:   email,
		Message: descriptionMessage,
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&deleteUserResponse)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"encoding of create user responce is failed",
			zap.Error(err),
		)
		return
	}
}

func New(deleteUserService DeleteUserService) *handler {
	return &handler{
		deleteUserService: deleteUserService,
	}
}

type handler struct {
	deleteUserService DeleteUserService
}
