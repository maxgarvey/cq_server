package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jonboulle/clockwork"
	"github.com/redis/go-redis/v9"
	"github.com/thanhpk/randstr"

	"github.com/maxgarvey/cq_server/config"
	"github.com/maxgarvey/cq_server/endpoints"
)

func main() {
	// Read in config.
	conf := config.GetConfig("localhost")

	// Connect to redis based off of config.
	redisClient := redis.NewClient(
		&redis.Options{
			Addr: fmt.Sprintf(
				"%s:%d",
				conf.Redis.Host,
				conf.Redis.Port,
			),
			Password: "", // no password set
			DB:       0,  // use default DB
		},
	)
	defer redisClient.Close()

	router := Router(clockwork.NewRealClock(), *redisClient)

	// Kick off endpoints.
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%d", conf.Server.Port),
			router))
}

// Router initialize router with endpoints.
func Router(clock clockwork.Clock, redisConnection redis.Client) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// Health check endpoint.
	router.HandleFunc("/health", endpoints.Health).Methods("GET")
	router.HandleFunc(
		"/ask/{requestType}",
		endpoints.Ask(clock, redisConnection, makeToken)).Methods("POST")
	router.HandleFunc(
		"/get/{id}",
		endpoints.Get(redisConnection)).Methods("GET")

	return router
}

func makeToken() string {
	return randstr.String(20)
}
