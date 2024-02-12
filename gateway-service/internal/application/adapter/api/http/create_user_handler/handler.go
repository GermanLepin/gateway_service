package create_user_handler

import (
	"context"
	"encoding/json"
	"net/http"

	"gateway-service/internal/application/dto"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	CretaeUserService interface {
		CreateUser(ctx context.Context, user *dto.User) error
	}

	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, status ...int) error
	}
)

func (h *handler) CretaeUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var cretaeUserRequest dto.CretaeUserRequest
	if err := json.NewDecoder(r.Body).Decode(&cretaeUserRequest); err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cretaeUserRequest.Password), 10)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	user := &dto.User{
		Name:     cretaeUserRequest.Name,
		Email:    cretaeUserRequest.Email,
		Password: string(passwordHash),
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
	err = encoder.Encode(&cretaeUserResponse)
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
