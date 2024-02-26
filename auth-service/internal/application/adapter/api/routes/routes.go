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

	LoginHandler interface {
		Login(w http.ResponseWriter, r *http.Request)
	}

	RefreshTokenHandler interface {
		RefreshToken(w http.ResponseWriter, r *http.Request)
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

	router.Route("/v1/api/user", func(r chi.Router) {
		//r.Post("/create", s.createUserHandler.CreateUser)
		// r.Post("/login", s.loginHandler.Login)
		// r.Post("/refresh-token", s.refreshTokenHandler.RefreshToken)
	})

	router.Route("/v1/api/user/protected", func(r chi.Router) {
		// r.Get("/fetch/{uuid}", s.fetchUserHandler.FetchUser)
		// r.Delete("/delete/{uuid}", s.deleteUserHandler.DeleteUser)
		// to do /logout
	})

	return router
}

func New(
	connection *sql.DB,

	// createUserHandler CreateUserHandler,
	// loginHandler LoginHandler,
	// refreshTokenHandler RefreshTokenHandler,
	// fetchUserHandler FetchUserHandler,
	// deleteUserHandler DeleteUserHandler,
) *service {
	return &service{
		// connection: connection,

		// createUserHandler:   createUserHandler,
		// loginHandler:        loginHandler,
		// refreshTokenHandler: refreshTokenHandler,
		// fetchUserHandler:    fetchUserHandler,
		// deleteUserHandler:   deleteUserHandler,
	}
}

type service struct {
	// 	connection *sql.DB

	// createUserHandler   CreateUserHandler
	// loginHandler        LoginHandler
	// refreshTokenHandler RefreshTokenHandler
	// fetchUserHandler    FetchUserHandler
	// deleteUserHandler   DeleteUserHandler
}
