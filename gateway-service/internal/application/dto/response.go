package dto

import "github.com/google/uuid"

type CretaeUserResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Name    string    `json:"name"`
	Message string    `json:"message"`
}

type PaymentResponse struct {
	OperationID uuid.UUID `json:"operation_id"`
	UserID      uuid.UUID `json:"user_id"`
	Status      string    `json:"status"`
	Error       string    `json:"error"`
}

type DeleteUserResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
