package json_service

import (
	"encoding/json"
	"log"
	"net/http"

	"client-service/internal/application/dto"
)

func (s *service) ErrorJSON(w http.ResponseWriter, paymentResponse dto.PaymentResponse, statusCode int) error {
	log.Printf("error: %s\n", paymentResponse.Error)

	out, err := json.Marshal(paymentResponse)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

type service struct{}

func New() *service {
	return &service{}
}
