package create_user_handler

import (
	"context"
	"encoding/json"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"

	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserService interface {
	CreateUser(ctx context.Context, user *dto.User) error
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	var createUserRequest dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("decoding of create user request is failed", zap.Error(err))
		return
	}

	logger = logger.With(
		zap.String("first_name", createUserRequest.FirstName),
		zap.String("last_name", createUserRequest.LastName),
		zap.String("password", createUserRequest.Password),
		zap.String("email", createUserRequest.Email),
		zap.Int("phone", createUserRequest.Phone),
		zap.String("user_type", createUserRequest.UserType),
	)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), 10)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("password hashing is failed", zap.Error(err))
		return
	}

	user := &dto.User{
		FirstName: createUserRequest.FirstName,
		LastName:  createUserRequest.LastName,
		Password:  string(passwordHash),
		Email:     createUserRequest.Email,
		Phone:     createUserRequest.Phone,
		UserType:  createUserRequest.UserType,
	}

	user.ID = uuid.New()
	if err := h.createUserService.CreateUser(ctx, user); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("user creation is failed", zap.Error(err))
		return
	}

	descriptionMessage := "user created successfully"
	createUserResponse := dto.CreateUserResponse{
		UserID:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		UserType:  user.UserType,
		Message:   descriptionMessage,
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&createUserResponse)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error("encoding of create user responce is failed", zap.Error(err))
		return
	}
}

func New(createUserService CreateUserService) *handler {
	return &handler{
		createUserService: createUserService,
	}
}

type handler struct {
	createUserService CreateUserService
}
