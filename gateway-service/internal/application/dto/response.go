package dto

import "github.com/google/uuid"

type CreateUserResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Phone   int       `json:"phone"`
	Email   string    `json:"email"`
	Message string    `json:"message"`
}

type FetchUserResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Phone   int       `json:"phone"`
	Email   string    `json:"email"`
	Message string    `json:"message"`
}

type LoginResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Phone    int       `json:"phone"`
	Email    string    `json:"email"`
	JWTToken string    `json:"jwt_token"`
	Message  string    `json:"message"`
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
