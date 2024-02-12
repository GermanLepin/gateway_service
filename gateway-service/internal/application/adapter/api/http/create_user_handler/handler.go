package create_user_handler

import (
	"context"
	"encoding/json"
	"net/http"

	"gateway-service/internal/application/dto"

	"github.com/google/uuid"
)

type (
	CretaeUserService interface {
		CreateUser(ctx context.Context, user dto.CretaeUserRequest) error
	}

	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, status ...int) error
	}
)

func (h *handler) CretaeUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var user dto.CretaeUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	user.ID = uuid.New()
	if err := h.cretaeUserService.CreateUser(ctx, user); err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	descriptionMessage := "user created successfully"
	cretaeUserResponse := dto.CretaeUserResponse{
		UserID:  user.ID,
		Name:    user.Name,
		Message: descriptionMessage,
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(&cretaeUserResponse)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
}

func New(
	cretaeUserService CretaeUserService,
	jsonService JsonService,
) *handler {
	return &handler{
		cretaeUserService: cretaeUserService,
		jsonService:       jsonService,
	}
}

type handler struct {
	cretaeUserService CretaeUserService
	jsonService       JsonService
}
