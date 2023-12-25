package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/benbjohnson/clock"
	"github.com/gorilla/mux"

	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/redis"
)

// Update updates the status of a task (called from workers)
func Update(
	clock clock.Clock, redisClient *redis.Redis, logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		PerformUpdate(
			w, r, clock, redisClient, logger,
		)
	}
}

// Make this a function so we can reuse it for AdminUpdate
func PerformUpdate(
	w http.ResponseWriter,
	r *http.Request,
	clock clock.Clock,
	redisClient *redis.Redis,
	logger *slog.Logger,
) {
	// Parse the requestType from URL.
	rawRequestType := mux.Vars(r)["requestType"]
	requestType := data.GetRequestType(rawRequestType)

	// Parse response id from URL.
	requestID := mux.Vars(r)["id"]

	// Read and unmarshal the incoming request.
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"error reading request body: %s\n",
				fmt.Errorf("%w", err),
			),
		)
		return
	}
	var request data.UpdateRequest
	json.Unmarshal(requestBody, &request)

	// Retrieve record from Redis
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

	// Update record of request.
	record.Status = data.Status(
		data.GetStatus(request.Status),
	)
	recordJSON, err := json.Marshal(record)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"error marshalling JSON for: %v\nerr: %s\n",
				record,
				fmt.Errorf("%w", err),
			),
		)
		return
	}
	ctx := context.Background()

	// Update redis.
	key := fmt.Sprintf(
		"%s:%s",
		requestType.String(),
		requestID,
	)
	err = redisClient.Set(
		ctx,
		key,
		recordJSON,
	)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"redis write failed for: %s\n%s\n%s\n",
				key,
				recordJSON,
				fmt.Errorf("%w", err),
			),
		)
		return
	}

	// Debug message.
	logger.Debug(
		fmt.Sprintf(
			"update endpoint requested. [requestType=%s][requestId=%s][status=%s]",
			requestType.String(),
			requestID,
			record.Status,
		),
	)

	json.NewEncoder(w).Encode(record.ToFullRecordResponse())
}
