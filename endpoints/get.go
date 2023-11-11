package endpoints

import (
	"context"
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
		// Parse response id from URL.
		requestID := mux.Vars(r)["id"]
		log.Printf(
			"requestID: %s",
			requestID,
		)

		// Retrieve response from DB.
		var response data.Record
		ctx := context.Background()
		response, err := redisClient.Get(
			ctx,
			fmt.Sprintf(
				"response:%s",
				requestID,
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(
			"get endpoint requested. [ID=%s, status=%s]",
			response.ID,
			fmt.Sprint(response.Status),
		)
		// If response is not ready.
		if response.Status != data.DONE {
			fmt.Fprintf(w, "not ready")
			return
		}

		// If response is ready, return it.
		fmt.Fprintf(w, response.Body)
	}
}
