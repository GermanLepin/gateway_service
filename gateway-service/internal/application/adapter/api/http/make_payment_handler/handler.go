package make_payment_handler

import (
	"time"

	"context"
	"encoding/json"
	"net/http"

	"gateway-service/internal/application/dto"

	"github.com/google/uuid"
)

type (
	MakePaymentService interface {
		Process(ctx context.Context, w http.ResponseWriter, paymentRequest dto.PaymentRequest) (string, error)
	}

	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, status ...int) error
	}
)

func (h *handler) MakePayment(w http.ResponseWriter, r *http.Request) {
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

	paymentStatus, err := h.makePaymentService.Process(ctx, w, paymentRequest)
	paymentResponse.Status = paymentStatus
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&paymentResponse)
	if err != nil {
		h.jsonService.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
}

func New(
	makePaymentService MakePaymentService,
	jsonService JsonService,
) *handler {
	return &handler{
		makePaymentService: makePaymentService,
		jsonService:        jsonService,
	}
}

type handler struct {
	makePaymentService MakePaymentService
	jsonService        JsonService
}
