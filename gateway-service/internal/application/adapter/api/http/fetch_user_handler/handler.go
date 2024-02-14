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

type (
	FetchUserService interface {
		FetchUser(ctx context.Context, userEmail string) (dto.User, error)
	}
)

func (h *handler) FetchUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	userEmail := chi.URLParam(r, "email")

	logger = logger.With(zap.String("email", userEmail))
	user, err := h.fetchUserService.FetchUser(ctx, userEmail)
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
		zap.String("name", user.Name),
		zap.String("surname", user.Surname),
		zap.Int("phone", user.Phone),
		zap.String("email", user.Email),
		zap.String("password", user.Password),
	)

	descriptionMessage := "user fetched successfully"
	deleteUserResponse := dto.FetchUserResponse{
		UserID:   user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: user.Password,
		Message:  descriptionMessage,
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
