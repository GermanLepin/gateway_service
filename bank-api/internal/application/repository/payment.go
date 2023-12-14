package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func (p *paymentRepository) SavePaymentInformation(ctx context.Context, operationID uuid.UUID, paymentStatus string) error {
	stmt := `insert into service.payment_status (operation_id, status) values ($1, $2);`

	err := p.db.QueryRowContext(ctx, stmt, operationID, paymentStatus)
	if err != nil {
		return err.Err()
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
