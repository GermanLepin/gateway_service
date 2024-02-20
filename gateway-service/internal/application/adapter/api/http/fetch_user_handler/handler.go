package fetch_user_handler

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

type FetchUserService interface {
	FetchUser(ctx context.Context, email string) (dto.User, error)
}

func (h *handler) FetchUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	email := chi.URLParam(r, "email")

	logger = logger.With(zap.String("email", email))
	user, err := h.fetchUserService.FetchUser(ctx, email)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"user fetching is failed",
			zap.Error(err),
		)
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

	descriptionMessage := "user fetched successfully"
	deleteUserResponse := dto.FetchUserResponse{
		UserID:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		UserType:  user.UserType,
		Message:   descriptionMessage,
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

func New(fetchUserService FetchUserService) *handler {
	return &handler{
		fetchUserService: fetchUserService,
	}
}

type handler struct {
	fetchUserService FetchUserService
}
