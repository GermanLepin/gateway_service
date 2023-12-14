package main

import (
	"bank-api/db/postgres/connection"
	"bank-api/internal/application/adapter/api/http/payment_handler"
	"bank-api/internal/application/adapter/api/routes"
	"bank-api/internal/application/repository"
	"bank-api/internal/application/service/json_service"
	"bank-api/internal/application/service/payment_service"

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
	payment_service := payment_service.New(payment_repository)
	payment_handler := payment_handler.New(payment_service, jsonService)

	api_routes := routes.New(connection, payment_handler)

	log.Printf("starting payment service on port %s\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: api_routes.NewRoutes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
