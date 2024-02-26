package main

import (
	"auth-service/db/postgres/connection"
	"auth-service/internal/application/adapter/api/routes"
	"auth-service/internal/application/helper/logging"

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
	//session_repository := repository.NewSessionRepository(connection)

	api_routes := routes.New(connection)

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
