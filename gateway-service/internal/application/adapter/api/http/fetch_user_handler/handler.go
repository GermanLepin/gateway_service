package fetch_user_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gateway-service/internal/application/dto"

	"github.com/go-chi/chi/v5"
)

type (
	FetchUserService interface {
		FetchUser(ctx context.Context, userEmail string) (dto.User, error)
	}

	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, status ...int) error
	}
)

func (h *handler) FetchUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userEmail := chi.URLParam(r, "email")
	user, err := h.fetchUserService.FetchUser(ctx, userEmail)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	descriptionMessage := "user fetched successfully"
	deleteUserResponse := dto.FetchUserResponse{
		UserID:  user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Phone:   user.Phone,
		Email:   user.Email,
		Message: descriptionMessage,
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&deleteUserResponse)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
}

func New(
	fetchUserService FetchUserService,
	jsonService JsonService,
) *handler {
	return &handler{
		fetchUserService: fetchUserService,
		jsonService:      jsonService,
	}
}

type handler struct {
	fetchUserService FetchUserService
	jsonService      JsonService
}
