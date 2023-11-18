package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/redis"
)

// Get a response.
func Get(redisClient *redis.Redis) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the requestType from URL.
		rawRequestType := mux.Vars(r)["requestType"]
		requestType := data.GetRequestType(rawRequestType)

		// Parse response id from URL.
		requestID := mux.Vars(r)["id"]
		log.Printf(
			"requestID: %s",
			requestID,
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
			log.Fatalf(
				"error retrieving record from redis: %s\n",
				fmt.Errorf("%w", err),
			)
		}

		log.Printf(
			"get endpoint requested. [requestType=%s, ID=%s, status=%s]",
			requestType.String(),
			record.ID,
			fmt.Sprint(record.Status),
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
