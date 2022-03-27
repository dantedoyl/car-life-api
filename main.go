package main

import (
	"github.com/dantedoyl/car-life-api/internal/app/clients/database"
	clubs_delivery "github.com/dantedoyl/car-life-api/internal/app/clubs/delivery/http"
	clubs_repository "github.com/dantedoyl/car-life-api/internal/app/clubs/repository/postgres"
	clubs_usecase "github.com/dantedoyl/car-life-api/internal/app/clubs/usecase"
	events_delivery "github.com/dantedoyl/car-life-api/internal/app/events/delivery/http"
	events_repository "github.com/dantedoyl/car-life-api/internal/app/events/repository/postgres"
	events_usecase "github.com/dantedoyl/car-life-api/internal/app/events/usecase"

	"github.com/dantedoyl/car-life-api/internal/app/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"time"

	_ "github.com/dantedoyl/car-life-api/docs"
	"github.com/gorilla/mux"
)

// @title           Swagger Example API
// @version         1.0
// @description     API for CarLife application

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	postgresDB, err := database.NewPostgres("host=localhost port=5432 user=postgres password=postgres dbname=car_life_api sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	eventsRepo := events_repository.NewProductRepository(postgresDB.GetDatabase())
	eventsUcase := events_usecase.NewEventsUsecase(eventsRepo)
	eventHandler := events_delivery.NewEventsHandler(eventsUcase)

	clubsRepo := clubs_repository.NewClubRepository(postgresDB.GetDatabase())
	clubsUcase := clubs_usecase.NewClubsUsecase(clubsRepo)
	clubsHandler := clubs_delivery.NewClubsHandler(clubsUcase)

	router := mux.NewRouter()

	static := router.PathPrefix("/static").Subrouter()
	static.Handle("/events/{key}", http.FileServer(http.Dir("."))).Methods(http.MethodGet)
	static.Handle("/clubs/{key}", http.FileServer(http.Dir("."))).Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.CorsControlMiddleware)
	eventHandler.Configure(api)
	clubsHandler.Configure(api)
	api.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
