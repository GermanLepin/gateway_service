package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type (
	UpdateStatusHandler interface {
		UpdateStatus(w http.ResponseWriter, r *http.Request)
	}

	PaymentHandler interface {
		Payment(w http.ResponseWriter, r *http.Request)
	}
)

func (s *service) NewRoutes() http.Handler {
	router := chi.NewRouter()

	// specify who is allowed to connect
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/", func(r chi.Router) {
		r.Post("/payment", s.paymentHandler.Payment)
	})

	return router
}

func New(
	connection *sql.DB,
	paymentHandler PaymentHandler,
) *service {
	return &service{
		connection:     connection,
		paymentHandler: paymentHandler,
	}
}

type service struct {
	connection     *sql.DB
	paymentHandler PaymentHandler
}
