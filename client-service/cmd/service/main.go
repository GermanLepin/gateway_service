package main

import (
	"client-service/db/postgres/connection"
	"client-service/internal/application/adapter/api/http/payment_handler"
	"client-service/internal/application/adapter/api/routes"
	"client-service/internal/application/repository"
	"client-service/internal/application/service/json_service"
	"client-service/internal/application/service/payment_service"

	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

const webPort = "80"

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file")
	}
}

func main() {
	connection := connection.StartDB()
	jsonService := json_service.New()

	payment_repository := repository.NewPaymentRepository(connection)
	make_payment_service := payment_service.New(payment_repository)
	make_payment_handler := payment_handler.New(make_payment_service, jsonService)

	api_routes := routes.New(connection, make_payment_handler)

	log.Printf("starting client service on port %s\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: api_routes.NewRoutes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
