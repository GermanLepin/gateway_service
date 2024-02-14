package make_payment_handler

import (
	"time"

	"context"
	"encoding/json"
	"net/http"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	MakePaymentService interface {
		Process(ctx context.Context, w http.ResponseWriter, paymentRequest dto.PaymentRequest) (string, error)
	}
)

func (h *handler) MakePayment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger := logging.LoggerFromContext(ctx)
	ctx = logging.ContextWithLogger(ctx, logger)

	var paymentRequest dto.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"decoding of payment request is failed",
			zap.Error(err),
		)
		return
	}

	paymentRequest.OperationID = uuid.New()
	paymentResponse := dto.PaymentResponse{
		OperationID: paymentRequest.OperationID,
		UserID:      paymentRequest.UserID,
	}

	paymentStatus, err := h.makePaymentService.Process(ctx, w, paymentRequest)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"payment is failed",
			zap.Error(err),
		)
		return
	}
	paymentResponse.Status = paymentStatus

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&paymentResponse)
	if err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"encoding of create user responce is failed",
			zap.Error(err),
		)
		return
	}
}

func New(makePaymentService MakePaymentService) *handler {
	return &handler{
		makePaymentService: makePaymentService,
	}
}

type handler struct {
	makePaymentService MakePaymentService
}
