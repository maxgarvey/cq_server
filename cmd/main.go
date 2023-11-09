package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jonboulle/clockwork"
	"github.com/thanhpk/randstr"

	"github.com/maxgarvey/cq_server/config"
	"github.com/maxgarvey/cq_server/endpoints"
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
)

func main() {
	// Read in config.
	conf := config.GetConfig(
		"localhost",
	)

	// Connect to redis based off of config.
	redisClient := redis.Init(
		conf.Redis,
	)
	// TODO: actually defer closing the connection.
	// defer redisClient.Close()

	// Connect to Rabbit MQ based off of config.
	var rabbitmqClient *rabbitmq.Rabbitmq
	var err error
	if conf.Rabbitmq.Host != "" && conf.Rabbitmq.Port != 0 {
		rabbitmqClient, err = rabbitmq.Init(
			conf.Rabbitmq,
		)
		if err != nil {
			log.Fatalf(
				"Error connecting to rabbit mq: %s",
				err.Error(),
			)
		}
	}

	router := Router(
		clockwork.NewRealClock(),
		rabbitmqClient,
		redisClient,
	)

	// Kick off endpoints.
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(
				":%d",
				conf.Server.Port,
			),
			router,
		),
	)
}

// Router initialize router with endpoints.
func Router(
	clock clockwork.Clock,
	rabbitClient *rabbitmq.Rabbitmq,
	redisClient *redis.Redis,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// Health check endpoint. Is the service running?
	router.HandleFunc(
		"/health",
		endpoints.Health,
	).Methods("GET")
	// Ask endpoint. Ask the service to do a job.
	router.HandleFunc(
		"/ask/{requestType}",
		endpoints.Ask(
			clock,
			*rabbitClient,
			redisClient,
			makeToken,
		),
	).Methods("POST")
	// Get endpoint. Check on a job.
	router.HandleFunc(
		"/get/{id}",
		endpoints.Get(
			redisClient,
		),
	).Methods("GET")

	return router
}

func makeToken() string {
	return randstr.String(20)
}
