package login_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gateway-service/internal/application/dto"
)

type (
	LoginService interface {
		Login(ctx context.Context, loginingUser *dto.User) (dto.User, error)
	}

	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, status ...int) error
	}
)

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	loginingUser := &dto.User{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}

	user, err := h.loginService.Login(ctx, loginingUser)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	descriptionMessage := "user logined successfully"
	loginResponse := dto.LoginResponse{
		UserID:   user.ID,
		Name:     user.Name,
		Email:    user.Email,
		JWTToken: user.JWTToken,
		Message:  descriptionMessage,
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&loginResponse)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
}

func New(
	loginService LoginService,
	jsonService JsonService,
) *handler {
	return &handler{
		loginService: loginService,
		jsonService:  jsonService,
	}
}

type handler struct {
	loginService LoginService
	jsonService  JsonService
}
