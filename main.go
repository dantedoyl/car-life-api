package main

import (
	"fmt"
	_ "github.com/dantedoyl/car-life-api/docs"
	"github.com/dantedoyl/car-life-api/internal/app/clients/database"
	"github.com/dantedoyl/car-life-api/internal/app/clients/vk"
	clubs_delivery "github.com/dantedoyl/car-life-api/internal/app/clubs/delivery/http"
	clubs_repository "github.com/dantedoyl/car-life-api/internal/app/clubs/repository/postgres"
	clubs_usecase "github.com/dantedoyl/car-life-api/internal/app/clubs/usecase"
	events_delivery "github.com/dantedoyl/car-life-api/internal/app/events/delivery/http"
	events_repository "github.com/dantedoyl/car-life-api/internal/app/events/repository/postgres"
	events_usecase "github.com/dantedoyl/car-life-api/internal/app/events/usecase"
	"github.com/dantedoyl/car-life-api/internal/app/middleware"
	mini_events_delivery "github.com/dantedoyl/car-life-api/internal/app/mini_events/delivery/http"
	mini_events_repository "github.com/dantedoyl/car-life-api/internal/app/mini_events/repository/postgres"
	mini_events_usecase "github.com/dantedoyl/car-life-api/internal/app/mini_events/usecase"
	users_delivery "github.com/dantedoyl/car-life-api/internal/app/users/delivery/http"
	users_repository "github.com/dantedoyl/car-life-api/internal/app/users/repository/postgres"
	users_usecase "github.com/dantedoyl/car-life-api/internal/app/users/usecase"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/tarantool/go-tarantool"
	"log"
	"net/http"
	"time"
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

	opts := tarantool.Opts{
		User: "admin",
		Pass: "pass",
	}
	tarConn, err := tarantool.Connect("127.0.0.1:3301", opts)
	if err != nil {
		fmt.Println("baa: Connection refused:", err)
		return
	}

	vkCl := vk.NewVKClient(
			"0918a4f20918a4f20918a4f27909633217009180918a4f26b34f910640ea9466b2c60dd",
		"1778a9046f1d83c8716dffe78116e0ce119cde5d96f7a2a8557299d0983bd349e91e13dea9cd920a2b5ed",
		)

	//________________________________
	//session map
	//var tarConn *tarantool.Connection = nil

	userRepo := users_repository.NewUserRepository(postgresDB.GetDatabase(), tarConn)
	userUcase := users_usecase.NewUsersUsecase(userRepo)
	userHandler := users_delivery.NewUserssHandler(userUcase)

	clubsRepo := clubs_repository.NewClubRepository(postgresDB.GetDatabase())
	clubsUcase := clubs_usecase.NewClubsUsecase(clubsRepo)
	clubsHandler := clubs_delivery.NewClubsHandler(clubsUcase, vkCl)

	eventsRepo := events_repository.NewProductRepository(postgresDB.GetDatabase())
	eventsUcase := events_usecase.NewEventsUsecase(eventsRepo)
	eventHandler := events_delivery.NewEventsHandler(eventsUcase, clubsUcase, vkCl)

	miniEventsRepo := mini_events_repository.NewMiniEventsRepository(postgresDB.GetDatabase())
	miniEventsUcase := mini_events_usecase.NewMiniEventsUsecase(miniEventsRepo)
	miniEventHandler := mini_events_delivery.NewMiniEventsHandler(miniEventsUcase)

	mw := middleware.NewMiddleware(userUcase)

	router := mux.NewRouter()

	static := router.PathPrefix("/img").Subrouter()
	static.Handle("/events/{key}", http.FileServer(http.Dir("."))).Methods(http.MethodGet)
	static.Handle("/clubs/{key}", http.FileServer(http.Dir("."))).Methods(http.MethodGet)
	static.Handle("/cars/{key}", http.FileServer(http.Dir("."))).Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.CorsControlMiddleware)
	eventHandler.Configure(api, mw)
	clubsHandler.Configure(api, mw)
	userHandler.Configure(api, mw)
	miniEventHandler.Configure(api, mw)
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
