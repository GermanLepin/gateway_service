package json_service

import (
	"encoding/json"
	"net/http"

	"gateway-service/internal/application/dto"
)

func (s *service) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload dto.JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return s.WriteJSON(w, statusCode, payload)
}

func (s *service) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
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
