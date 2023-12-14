package payment_service

import (
	"context"
	"fmt"
	"net/http"

	"bank-api/internal/application/dto"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	SavePaymentInformation(ctx context.Context, operationID uuid.UUID, paymentStatus string) error
}

func (s *service) Process(ctx context.Context, w http.ResponseWriter, paymentRequest dto.PaymentRequest) (string, error) {
	paymentStatus := "succeed"
	err := s.paymentRepository.SavePaymentInformation(ctx, paymentRequest.OperationID, paymentStatus)
	if err != nil {
		return "error", fmt.Errorf("cannot save the data: %s", err)
	}

	return paymentStatus, nil
}

func New(paymentRepository PaymentRepository) *service {
	return &service{
		paymentRepository: paymentRepository,
	}
}

type service struct {
	paymentRepository PaymentRepository
}
