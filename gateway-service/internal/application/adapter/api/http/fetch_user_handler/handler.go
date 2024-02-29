package fetch_user_handler

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

type FetchUserService interface {
	FetchUser(ctx context.Context, r *http.Request, userId uuid.UUID) (dto.User, error)
}

func (h *handler) FetchUser(w http.ResponseWriter, r *http.Request) {
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
	user, err := h.fetchUserService.FetchUser(ctx, r, userUUID)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("user fetching is failed", zap.Error(err))
		return
	}

	logger = logger.With(
		zap.String("userID", user.ID.String()),
		zap.String("first_name", user.FirstName),
		zap.String("last_name", user.LastName),
		zap.String("password", user.Password),
		zap.String("email", user.Email),
		zap.Int("phone", user.Phone),
		zap.String("user_type", user.UserType),
	)

	fetchUserResponse := dto.FetchUserResponse{
		UserID:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		UserType:  user.UserType,
	}

	if err = jsonwrapper.WriteJSON(w, http.StatusOK, fetchUserResponse); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("could not send a fetch user response", zap.Error(err))
		return
	}
}

func New(fetchUserService FetchUserService) *handler {
	return &handler{
		fetchUserService: fetchUserService,
	}
}

type handler struct {
	fetchUserService FetchUserService
}
