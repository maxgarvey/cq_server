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
		logger.Debug(
			fmt.Sprintf(
				"requestID: %s",
				requestID,
			),
		)

		// Retrieve record from DB.
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
			logger.Error(
				fmt.Sprintf(
					"error retrieving record from redis: %s\n",
					fmt.Errorf("%w", err),
				),
			)
		}

		logger.Debug(
			fmt.Sprintf(
				"get endpoint requested. [requestType=%s, ID=%s, status=%s]",
				requestType.String(),
				record.ID,
				fmt.Sprint(record.Status),
			),
		)
		// If response is not ready.
		if record.Status != data.DONE {
			fmt.Fprintf(w, "not ready")
			return
		}

		// If response is ready, return it.
		json.NewEncoder(w).Encode(record.ToGetResponse())
	}
}
