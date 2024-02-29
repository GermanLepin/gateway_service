package delete_user_handler

import (
	"context"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DeleteUserService interface {
	DeleteUser(ctx context.Context, userId uuid.UUID) error
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	userId := chi.URLParam(r, "user_id")
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("user id parsing is failing", zap.Error(err))
		return
	}

	logger = logger.With(zap.String("uuid", userId))
	err = h.deleteUserService.DeleteUser(ctx, userUUID)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("user deletion is failed", zap.Error(err))
		return
	}

	descriptionMessage := "user deleted successfully"
	deleteUserResponse := dto.DeleteUserResponse{
		UserID:  userUUID,
		Message: descriptionMessage,
	}

	if err = jsonwrapper.WriteJSON(w, http.StatusOK, deleteUserResponse); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("could not send a delete user response.", zap.Error(err))
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
