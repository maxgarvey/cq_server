package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/benbjohnson/clock"
	"github.com/gorilla/mux"

	"github.com/maxgarvey/cq_server/admin"
	"github.com/maxgarvey/cq_server/config"
	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/endpoints"
	"github.com/maxgarvey/cq_server/postgres"
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
	"github.com/maxgarvey/cq_server/worker"
)

func main() {
	// Accept flag for config file.
	configFile := flag.String(
		"config",
		"./config/example.yaml",
		"location of yaml configuration file.",
	)
	flag.Parse()

	// Read in config.
	conf := config.GetConfig(
		*configFile,
	)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	clock := clock.New()

	// Connect to postgres based off of config.
	postgresClient := postgres.ConfigInit(
		conf.Postgres,
		clock,
		logger,
	)

	// Connect to redis based off of config.
	redisClient := redis.Init(
		conf.Redis,
	)
	defer redisClient.Close()

	// Connect to Rabbit MQ based off of config.
	var rabbitmqClient *rabbitmq.Rabbitmq
	var err error
	if conf.Rabbitmq.Host != "" && conf.Rabbitmq.Port != 0 {
		rabbitmqClient, err = rabbitmq.Init(
			conf.Rabbitmq,
		)
		if err != nil {
			logger.Error(
				fmt.Sprintf(
					"Error connecting to rabbit mq: %s\n",
					err.Error(),
				),
			)
		}
		defer rabbitmqClient.Close()
	}

	admin := admin.Admin{
		Clock:    clock,
		Postgres: &postgresClient,
		Logger:   logger,
	}

	// Initialize router to handle web requests.
	router := Router(
		clock,
		rabbitmqClient,
		redisClient,
		&admin,
		logger,
	)

	// Initialize queue worker.
	queueWorker := worker.Init(rabbitmqClient, redisClient, logger)
	// Start work.
	// TODO: add multithreading for workers
	queueWorker.Work()

	logger.Info(
		fmt.Sprintf(
			"Server listening on 127.0.0.1:%d\n",
			conf.Server.Port,
		),
	)

	// Kick off endpoints.
	err = http.ListenAndServe(
		fmt.Sprintf(
			":%d",
			conf.Server.Port,
		),
		router,
	)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error running server: %s\n",
				err.Error(),
			),
		)
	}
}

// Router initialize router with endpoints.
func Router(
	clock clock.Clock,
	rabbitClient *rabbitmq.Rabbitmq,
	redisClient *redis.Redis,
	admin *admin.Admin,
	logger *slog.Logger,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// Health check endpoint. Is the service running?
	router.HandleFunc(
		"/health",
		endpoints.Health(logger),
	).Methods("GET")
	// Ask endpoint. Ask the service to do a job.
	router.HandleFunc(
		"/ask/{requestType}",
		endpoints.Ask(
			clock,
			*rabbitClient,
			redisClient,
			data.MakeToken,
			logger,
		),
	).Methods("POST")
	// Get endpoint. Check on a job.
	router.HandleFunc(
		"/get/{requestType}/{id}",
		endpoints.Get(
			redisClient,
			logger,
		),
	).Methods("GET")

	// Admin endpoints
	router.HandleFunc(
		"/admin/login",
		endpoints.AdminLogin(
			admin,
			*logger,
		),
	).Methods("POST")
	router.HandleFunc(
		"/admin/get/{requestType}/{id}",
		endpoints.AdminGet(
			admin,
			*redisClient,
			*logger,
		),
	).Methods("GET")
	return router
}
