package dto

import "github.com/google/uuid"

type PaymentResponse struct {
	OperationID uuid.UUID `json:"operation_id"`
	UserID      uuid.UUID `json:"user_id"`
	Status      string    `json:"status"`
	Error       string    `json:"error"`
}
