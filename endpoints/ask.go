package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jonboulle/clockwork"

	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
)

// Ask enqueues a request and creates an entry in redis to track it.
func Ask(clock clockwork.Clock, rabbitmq rabbitmq.Rabbit, redisClient *redis.Redis, token func() string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestType := mux.Vars(r)["requestType"]
		token := token()

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Create redis record of request.
		response := &data.Record{
			Body:        string(requestBody),
			ID:          token,
			RequestType: requestType,
			Status:      data.IN_PROGRESS,
			Timestamp:   clock.Now().Unix(),
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		// Put it into redis.
		redisClient.Set(
			ctx,
			fmt.Sprintf(
				"%s:%s",
				requestType,
				token,
			),
			responseJSON,
		)

		// enqueue message to perform the work
		rabbitmq.Publish(string(responseJSON))

		log.Printf(
			"ask endpoint requested. [requestType=%s]",
			requestType,
		)

		// Return token associated with this request.
		askResp := &data.AskResponse{
			ID: token,
		}
		json.NewEncoder(w).Encode(&askResp)
	}
}
