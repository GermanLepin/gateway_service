package update_status_handler

import (
	"time"

	"context"
	"encoding/json"
	"net/http"

	"gateway-service/internal/application/dto"
	"gateway-service/internal/application/helper/jsonwrapper"
	"gateway-service/internal/application/helper/logging"

	"go.uber.org/zap"
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

	if err := h.paymentService.Process(ctx, paymentRequest); err != nil {
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
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
		jsonwrapper.ErrorJSON(w, err, http.StatusInternalServerError)
		logger.Error(
			"encoding of create user responce is failed",
			zap.Error(err),
		)
		return
	}
}

func New(paymentService PaymentService) *handler {
	return &handler{
		paymentService: paymentService,
	}
}

type handler struct {
	paymentService PaymentService
}
