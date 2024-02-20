package dto

import (
	"github.com/google/uuid"
)

type CreateUserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     int       `json:"phone"`
	UserType  string    `json:"user_type"`
	Message   string    `json:"message"`
}

type FetchUserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     int       `json:"phone"`
	UserType  string    `json:"user_type"`
	Message   string    `json:"message"`
}

type LoginResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     int       `json:"phone"`
	UserType  string    `json:"user_type"`
	Message   string    `json:"message"`
}

type DeleteUserResponse struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// TODO payment service
type PaymentResponse struct {
	OperationID uuid.UUID `json:"operation_id"`
	UserID      uuid.UUID `json:"user_id"`
	Status      string    `json:"status"`
	Error       string    `json:"error"`
}
