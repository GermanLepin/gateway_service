package main

import (
	"context"
	"fmt"

	"gateway-service/db/postgres/connection"
	"gateway-service/internal/application/adapter/api/http/create_user_handler"
	"gateway-service/internal/application/adapter/api/http/delete_user_handler"
	"gateway-service/internal/application/adapter/api/http/fetch_user_handler"
	"gateway-service/internal/application/adapter/api/http/login_handler"
	"gateway-service/internal/application/adapter/api/routes"
	"gateway-service/internal/application/helper/logging"
	"gateway-service/internal/application/repository"
	"gateway-service/internal/application/service/create_session_service"
	"gateway-service/internal/application/service/create_user_service"
	"gateway-service/internal/application/service/delete_user_service"
	"gateway-service/internal/application/service/fetch_user_service"
	"gateway-service/internal/application/service/login_service"

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

	user_repository := repository.NewUserRepository(connection)
	session_repository := repository.NewSessionRepository(connection)

	create_user_service := create_user_service.New(user_repository)
	create_session_service := create_session_service.New(session_repository)
	login_service := login_service.New(user_repository, create_session_service)
	fetch_user_service := fetch_user_service.New(user_repository)
	delete_user_service := delete_user_service.New(user_repository)

	create_user_handler := create_user_handler.New(create_user_service)
	login_handler := login_handler.New(login_service)
	fetch_user_handler := fetch_user_handler.New(fetch_user_service)
	delete_user_handler := delete_user_handler.New(delete_user_service)

	api_routes := routes.New(
		connection,
		create_user_handler,
		login_handler,
		fetch_user_handler,
		delete_user_handler,
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
