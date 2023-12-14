package dto

import "github.com/google/uuid"

type PaymentRequest struct {
	OperationID    uuid.UUID `json:"operation_id"`
	UserID         uuid.UUID `json:"user_id"`
	Amount         float32   `json:"amount"`
	CardNumber     int       `json:"card_number"`
	CardHolderName string    `json:"card_holder_number"`
	CVV            int       `json:"cvv"`
}
