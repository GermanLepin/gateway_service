package repository

import (
	"context"
	"database/sql"

	"gateway-service/internal/application/dto"

	"github.com/google/uuid"
)

func (p *paymentRepository) SavePaymentInformation(ctx context.Context, paymentRequest dto.PaymentRequest) error {
	stmt := `insert into service.payment_information (operation_id, user_id, amount, card_number, card_holder_name, cvv, payment_status) 
			values ($1, $2, $3, $4, $5, $6, $7);`

	err := p.db.QueryRowContext(ctx, stmt,
		paymentRequest.OperationID,
		paymentRequest.UserID,
		paymentRequest.Amount,
		paymentRequest.CardNumber,
		paymentRequest.CardHolderName,
		paymentRequest.CVV,
		paymentRequest.PaymentStatus,
	)
	if err != nil {
		return err.Err()
	}

	return nil
}

func (p *paymentRepository) UpdatePaymentStatus(ctx context.Context, operationID uuid.UUID, paymentStatus string) error {
	stmt := `update service.payment_information set payment_status = $1 where operation_id = $2;`

	if _, err := p.db.ExecContext(ctx, stmt, paymentStatus, operationID); err != nil {
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
