package update_status_handler

import (
	"time"

	"context"
	"encoding/json"
	"net/http"

	"gateway-service/internal/application/dto"
)

type (
	PaymentService interface {
		Process(ctx context.Context, paymentRequest dto.PaymentRequest) error
	}

	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, statusCode int) error
	}
)

func (h *handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var paymentRequest dto.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err := h.paymentService.Process(ctx, paymentRequest); err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	paymentResponse := dto.PaymentResponse{
		OperationID: paymentRequest.OperationID,
		UserID:      paymentRequest.UserID,
		Status:      "",
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(&paymentResponse)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
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
