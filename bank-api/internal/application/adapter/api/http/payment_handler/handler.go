package payment_handler

import (
	"time"

	"context"
	"encoding/json"
	"net/http"

	"bank-api/internal/application/dto"

	"github.com/google/uuid"
)

type (
	PaymentService interface {
		Process(ctx context.Context, w http.ResponseWriter, paymentRequest dto.PaymentRequest) (string, error)
	}

	JsonService interface {
		ErrorJSON(w http.ResponseWriter, paymentResponse dto.PaymentResponse, statusCode int) error
	}
)

func (h *handler) Payment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var paymentRequest dto.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	paymentRequest.OperationID = uuid.New()

	paymentResponse := dto.PaymentResponse{
		OperationID: paymentRequest.OperationID,
		UserID:      paymentRequest.UserID,
	}

	paymentStatus, err := h.paymentService.Process(ctx, w, paymentRequest)
	paymentResponse.Status = paymentStatus
	if err != nil {
		paymentResponse.Error = err.Error()

		h.jsonService.ErrorJSON(w, paymentResponse, http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&paymentResponse)
	if err != nil {
		paymentResponse.Error = err.Error()

		h.jsonService.ErrorJSON(w, paymentResponse, http.StatusInternalServerError)
		return
	}
}

func New(
	paymentService PaymentService,
	jsonService JsonService,
) *handler {
	return &handler{
		paymentService: paymentService,
		jsonService:    jsonService,
	}
}

type handler struct {
	paymentService PaymentService
	jsonService    JsonService
}
