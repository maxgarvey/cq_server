package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/redis"
)

// Get a response.
func Get(redisClient *redis.Redis, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the requestType from URL.
		rawRequestType := mux.Vars(r)["requestType"]
		requestType := data.GetRequestType(rawRequestType)

		// Parse response id from URL.
		requestID := mux.Vars(r)["id"]

		record, err := GetFromRedis(requestID, requestType, *redisClient)
		if err != nil {
			logger.Error(
				fmt.Sprintf(
					"error retrieving record from redis: %s\n",
					fmt.Errorf("%w", err),
				),
			)
			return
		}

		// If response is ready, return it.
		json.NewEncoder(w).Encode(record.ToGetResponse())
	}
}

func GetFromRedis(
	requestID string, requestType data.RequestType, redisClient redis.Redis,
) (data.Record, error) {
	// Retrieve record from Redis.
	var record data.Record
	ctx := context.Background()
	record, err := redisClient.Get(
		ctx,
		fmt.Sprintf(
			"%s:%s",
			requestType.String(),
			requestID,
		),
	)
	if err != nil {
		return record, err
	}

	return record, nil
}
