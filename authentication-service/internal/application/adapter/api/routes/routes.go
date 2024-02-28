package routes

import (
	"database/sql"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type (
	LoginHandler interface {
		Login(w http.ResponseWriter, r *http.Request)
	}

	ValidateTokenHandler interface {
		ValidateToken(w http.ResponseWriter, r *http.Request)
	}

	RefreshTokenHandler interface {
		RefreshToken(w http.ResponseWriter, r *http.Request)
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
		r.Post("/login", s.loginHandler.Login)
		r.Post("/validate-token", s.validateTokenHandler.ValidateToken)
		r.Post("/refresh-token", s.refreshTokenHandler.RefreshToken)
	})

	return router
}

func New(
	connection *sql.DB,

	loginHandler LoginHandler,
	validateTokenHandler ValidateTokenHandler,
	refreshTokenHandler RefreshTokenHandler,
) *service {
	return &service{
		connection: connection,

		loginHandler:         loginHandler,
		validateTokenHandler: validateTokenHandler,
		refreshTokenHandler:  refreshTokenHandler,
	}
}

type service struct {
	connection *sql.DB

	loginHandler         LoginHandler
	validateTokenHandler ValidateTokenHandler
	refreshTokenHandler  RefreshTokenHandler
}
