package main

import (
	"github.com/dantedoyl/car-life-api/internal/app/clients/database"
	delivery "github.com/dantedoyl/car-life-api/internal/app/events/delivery/http"
	events_repository "github.com/dantedoyl/car-life-api/internal/app/events/repository/postgres"
	"github.com/dantedoyl/car-life-api/internal/app/events/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	postgresDB, err := database.NewPostgres("host=localhost port=5432 user=postgres password=ysnpkoyapassword dbname=car-life-api sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	eventsRepo := events_repository.NewProductRepository(nil)
	eventsUcase := usecase.NewEventsUsecase(eventsRepo)
	eventHandler := delivery.NewEventsHandler(eventsUcase)

	router := mux.NewRouter()



	api := router.PathPrefix("/api/v1").Subrouter()
	eventHandler.Configure(api)

	router.HandleFunc("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/swagger.json"), //The url pointing to API definition
	)).Methods(http.MethodGet, http.MethodOptions)


	server := http.Server{
		Addr:         ":8080",
		Handler:      api,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
