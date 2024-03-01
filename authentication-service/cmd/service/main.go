package main

import (
	"authentication-service/db/postgres/connection"
	"authentication-service/internal/application/adapter/api/http/login_user_handler"
	"authentication-service/internal/application/adapter/api/http/refresh_token_handler"
	"authentication-service/internal/application/adapter/api/http/validate_token_handler"
	"authentication-service/internal/application/adapter/api/routes"
	"authentication-service/internal/application/helper/logging"
	"authentication-service/internal/application/repository"
	"authentication-service/internal/application/service/create_session_service"
	"authentication-service/internal/application/service/refresh_token_service"
	"authentication-service/internal/application/service/validate_token_service"

	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

const webPort = "80"

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file")
	}
}

func main() {
	_, ctx := errgroup.WithContext(context.Background())
	_, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	defer stop()

	_, loggerSyncFunc, err := logging.ZapFromEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer loggerSyncFunc()

	connection := connection.StartDB()
	session_repository := repository.NewSessionRepository(connection)

	create_session_service := create_session_service.New(session_repository)
	validate_token_service := validate_token_service.New(session_repository)
	refresh_token_service := refresh_token_service.New(session_repository, create_session_service)

	login_user_handler := login_user_handler.New(create_session_service)
	validate_token_handler := validate_token_handler.New(validate_token_service)
	refresh_token_handler := refresh_token_handler.New(refresh_token_service)

	api_routes := routes.New(
		connection,
		login_user_handler,
		validate_token_handler,
		refresh_token_handler,
	)

	log.Printf("starting client service on port %s\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: api_routes.NewRoutes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
