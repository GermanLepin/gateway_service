package repository

import (
	"context"
	"database/sql"

	"payment-service/internal/application/dto"

	"github.com/google/uuid"
)

func (p *paymentRepository) SavePaymentInformation(ctx context.Context, paymentRequest dto.PaymentRequest) error {
	stmt := `insert into service.actual_payment_information (operation_id, user_id, amount, card_number, card_holder_name, cvv, payment_status, operation_status) 
			values ($1, $2, $3, $4, $5, $6, $7, $8);`

	err := p.db.QueryRowContext(ctx, stmt,
		paymentRequest.OperationID,
		paymentRequest.UserID,
		paymentRequest.Amount,
		paymentRequest.CardNumber,
		paymentRequest.CardHolderName,
		paymentRequest.CVV,
		paymentRequest.PaymentStatus,
		paymentRequest.OperationStatus,
	)
	if err != nil {
		return err.Err()
	}

	return nil
}

func (p *paymentRepository) UpdateStatuses(ctx context.Context, operationID uuid.UUID, paymentStatus, operationStatus string) error {
	stmt := `update service.actual_payment_information set 
		payment_status = $1,  
		operation_status = $2
		where operation_id = $3;`

	_, err := p.db.ExecContext(ctx, stmt,
		paymentStatus,
		operationStatus,
		operationID,
	)
	if err != nil {
		return err
	}

	return nil
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *paymentRepository {
	return &paymentRepository{
		db: db,
	}
}
