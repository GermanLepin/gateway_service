package dto

import (
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     int    `json:"phone"`
	UserType  string `json:"user_type"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO payment service
type PaymentRequest struct {
	OperationID    uuid.UUID `json:"operation_id"`
	UserID         uuid.UUID `json:"user_id"`
	Amount         float32   `json:"amount"`
	CardNumber     uint32    `json:"card_number"`
	CVV            uint32    `json:"cvv"`
	CardHolderName string    `json:"card_holder_name"`
	PaymentStatus  string    `json:"payment_status"`
}
