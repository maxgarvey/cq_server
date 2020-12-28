package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jonboulle/clockwork"
	"github.com/thanhpk/randstr"

	"github.com/maxgarvey/cq_server/config"
	"github.com/maxgarvey/cq_server/endpoints"
)

func main() {
	// Read in config.
	conf := config.GetConfig("localhost")

	// Connect to redis based off of config.
	redisConnection, err := redis.Dial(
		conf.Redis.ConnectionType,
		fmt.Sprintf(
			"%s:%d", conf.Redis.Host, conf.Redis.Port))
	if err != nil {
		log.Printf("Error connecting to redis. [err=%v]", err)
		return
	}
	defer redisConnection.Close()

	router := Router(clockwork.NewRealClock(), redisConnection)

	// Kick off endpoints.
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%d", conf.Server.Port),
			router))
}

// Router initialize router with endpoints.
func Router(clock clockwork.Clock, redisConnection redis.Conn) *mux.Router {
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
