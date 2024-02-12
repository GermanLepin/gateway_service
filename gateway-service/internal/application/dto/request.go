package dto

import "github.com/google/uuid"

type CretaeUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PaymentRequest struct {
	OperationID    uuid.UUID `json:"operation_id"`
	UserID         uuid.UUID `json:"user_id"`
	Amount         float32   `json:"amount"`
	CardNumber     uint32    `json:"card_number"`
	CVV            uint32    `json:"cvv"`
	CardHolderName string    `json:"card_holder_name"`
	PaymentStatus  string    `json:"payment_status"`
}
