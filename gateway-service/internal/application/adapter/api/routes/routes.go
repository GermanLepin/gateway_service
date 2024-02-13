package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type (
	CreateUserHandler interface {
		CreateUser(w http.ResponseWriter, r *http.Request)
	}

	FetchUserHandler interface {
		FetchUser(w http.ResponseWriter, r *http.Request)
	}

	LoginHandler interface {
		Login(w http.ResponseWriter, r *http.Request)
	}

	MakePaymentHandler interface {
		MakePayment(w http.ResponseWriter, r *http.Request)
	}

	UpdateStatusHandler interface {
		UpdateStatus(w http.ResponseWriter, r *http.Request)
	}

	DeleteUserHandler interface {
		DeleteUser(w http.ResponseWriter, r *http.Request)
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

	router.Route("/user", func(r chi.Router) {
		r.Post("/create", s.createUserHandler.CreateUser)
		r.Post("/login", s.loginHandler.Login)
		r.Get("/fetch/{string}", s.fetchUserHandler.FetchUser)
		r.Delete("/delete/{uuid}", s.deleteUserHandler.DeleteUser)
	})

	router.Route("/", func(r chi.Router) {
		r.Post("/payment", s.paymentHandler.MakePayment)
	})

	return router
}

func New(
	connection *sql.DB,

	createUserHandler CreateUserHandler,
	fetchUserHandler FetchUserHandler,
	loginHandler LoginHandler,
	makePaymentHandler MakePaymentHandler,
	deleteUserHandler DeleteUserHandler,
) *service {
	return &service{
		connection: connection,

		createUserHandler: createUserHandler,
		fetchUserHandler:  fetchUserHandler,
		loginHandler:      loginHandler,
		paymentHandler:    makePaymentHandler,
		deleteUserHandler: deleteUserHandler,
	}
}

type service struct {
	connection *sql.DB

	createUserHandler CreateUserHandler
	fetchUserHandler  FetchUserHandler
	loginHandler      LoginHandler
	paymentHandler    MakePaymentHandler
	deleteUserHandler DeleteUserHandler
}
