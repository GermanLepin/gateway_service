package payment_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"client-service/internal/application/dto"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	SavePaymentInformation(ctx context.Context, paymentRequest dto.PaymentRequest) error
	UpdatePaymentStatus(ctx context.Context, operationID uuid.UUID, paymentStatus string) error
}

func (s *service) Process(ctx context.Context, w http.ResponseWriter, paymentRequest dto.PaymentRequest) (status string, err error) {
	paymentRequest.PaymentStatus = "processing"
	err = s.paymentRepository.SavePaymentInformation(ctx, paymentRequest)
	if err != nil {
		return "error", fmt.Errorf("cannot save the payment information: %s", err)
	}

	paymenStatus, errorPaymentRequest := s.paymentRequest(ctx, w, paymentRequest)
	err = s.paymentRepository.UpdatePaymentStatus(ctx, paymentRequest.OperationID, paymenStatus)
	if err != nil {
		return "error", fmt.Errorf("cannot update the status: %s", err)
	}

	return paymenStatus, errorPaymentRequest
}

func (s *service) paymentRequest(ctx context.Context, w http.ResponseWriter, paymentRequest dto.PaymentRequest) (status string, err error) {
	paymentServiceURL := "http://payment-service/payment"
	jsonData, err := json.MarshalIndent(paymentRequest, "", "\t")
	if err != nil {
		return "error", fmt.Errorf("client: could not create request: %s", err)
	}

	request, err := http.NewRequest(http.MethodPost, paymentServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "error", fmt.Errorf("client: could not create request: %s", err)
	}
	request.Close = true

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "error", fmt.Errorf("client: error making http request: %s", err)
	}
	defer response.Body.Close()

	var paymentResponse dto.PaymentResponse
	err = json.NewDecoder(response.Body).Decode(&paymentResponse)
	if err != nil {
		return "error", fmt.Errorf("client: could not read response body: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		return "error", fmt.Errorf(paymentResponse.Error)
	}

	return paymentResponse.Status, nil
}

func New(paymentRepository PaymentRepository) *service {
	return &service{
		paymentRepository: paymentRepository,
	}
}

type service struct {
	paymentRepository PaymentRepository
}
