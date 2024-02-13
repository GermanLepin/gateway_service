package create_user_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/logging"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type (
	CreateUserService interface {
		CreateUser(ctx context.Context, user *dto.User) error
	}

	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, status ...int) error
	}
)

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	var createUserRequest dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"decoding of create user request is failed",
			zap.Error(err),
		)
		return
	}

	logger = logger.With(
		zap.String("name", createUserRequest.Name),
		zap.String("surname", createUserRequest.Surname),
		zap.Int("phone", createUserRequest.Phone),
		zap.String("email", createUserRequest.Email),
		zap.String("password", createUserRequest.Password),
	)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), 10)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"password hashing is failed",
			zap.Error(err),
		)
		return
	}

	user := &dto.User{
		Name:     createUserRequest.Name,
		Surname:  createUserRequest.Surname,
		Phone:    createUserRequest.Phone,
		Email:    createUserRequest.Email,
		Password: string(passwordHash),
	}

	user.ID = uuid.New()
	if err := h.createUserService.CreateUser(ctx, user); err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"user creation is failed",
			zap.Error(err),
		)
		return
	}

	descriptionMessage := "user created successfully"
	createUserResponse := dto.CreateUserResponse{
		UserID:  user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Phone:   user.Phone,
		Email:   user.Email,
		Message: descriptionMessage,
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&createUserResponse)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"encoding of create user responce is failed",
			zap.Error(err),
		)
		return
	}
}

func New(
	createUserService CreateUserService,
	jsonService JsonService,
) *handler {
	return &handler{
		createUserService: createUserService,
		jsonService:       jsonService,
	}
}

type handler struct {
	createUserService CreateUserService
	jsonService       JsonService
}
