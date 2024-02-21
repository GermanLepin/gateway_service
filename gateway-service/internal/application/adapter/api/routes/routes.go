package routes

import (
	"database/sql"
	middleware "gateway-service/gateway-service/internal/application/adapter/api/middleware/validate_jwt_token"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type (
	CreateUserHandler interface {
		CreateUser(w http.ResponseWriter, r *http.Request)
	}

	LoginHandler interface {
		Login(w http.ResponseWriter, r *http.Request)
	}

	FetchUserHandler interface {
		FetchUser(w http.ResponseWriter, r *http.Request)
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

	router.Route("/v1/user", func(r chi.Router) {
		r.Post("/create", s.createUserHandler.CreateUser)
		r.Post("/login", s.loginHandler.Login)
	})

	router.Route("/v1/user/protected", func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Get("/fetch/{email}", s.fetchUserHandler.FetchUser)
		r.Delete("/delete/{email}", s.deleteUserHandler.DeleteUser)
	})

	return router
}

func New(
	connection *sql.DB,

	createUserHandler CreateUserHandler,
	loginHandler LoginHandler,
	fetchUserHandler FetchUserHandler,
	deleteUserHandler DeleteUserHandler,
) *service {
	return &service{
		connection: connection,

		createUserHandler: createUserHandler,
		loginHandler:      loginHandler,
		fetchUserHandler:  fetchUserHandler,
		deleteUserHandler: deleteUserHandler,
	}
}

type service struct {
	connection *sql.DB

	createUserHandler CreateUserHandler
	loginHandler      LoginHandler
	fetchUserHandler  FetchUserHandler
	deleteUserHandler DeleteUserHandler
}
