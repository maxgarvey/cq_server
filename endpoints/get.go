package endpoints

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"github.com/maxgarvey/cq_server/data"
)

// Get a response.
func Get(redisClient redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse response id from URL.
		requestID := mux.Vars(r)["id"]
		log.Printf("requestID: %s", requestID)

		// Retrieve raw response from DB.
		var response data.Response
		ctx := context.Background()
		err := redisClient.Get(
			ctx,
			fmt.Sprintf("response:%s", requestID),

		// Parse response.
		).Scan(&response)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(
			"get endpoint requested. [ID=%s, status=%s]",
			response.ID,
			response.Status)
		// If response is not ready.
		if response.Status != "DONE" {
			fmt.Fprintf(w, "not ready")
			return
		}

		// If response is ready, return it.
		fmt.Fprintf(w, response.Body)
	}
}
