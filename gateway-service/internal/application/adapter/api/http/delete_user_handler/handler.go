package delete_user_handler

import (
	"context"
	"encoding/json"
	"net/http"

	"gateway-service/internal/application/dto"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type (
	DeleteUserService interface {
		DeleteUser(ctx context.Context, userID uuid.UUID) error
	}

	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, status ...int) error
	}
)

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	userID := chi.URLParam(r, "uuid")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = h.deleteUserService.DeleteUser(ctx, userUUID)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	descriptionMessage := "user deleted successfully"
	deleteUserResponse := dto.DeleteUserResponse{
		UserID:  userUUID,
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
	deleteUserService DeleteUserService,
	jsonService JsonService,
) *handler {
	return &handler{
		deleteUserService: deleteUserService,
		jsonService:       jsonService,
	}
}

type handler struct {
	deleteUserService DeleteUserService
	jsonService       JsonService
}
