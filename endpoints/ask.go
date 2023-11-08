package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jonboulle/clockwork"
	"github.com/redis/go-redis/v9"

	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/rabbitmq"
)

// Ask enqueues a request and creates an entry in redis to track it.
func Ask(clock clockwork.Clock, rabbitmq rabbitmq.Rabbitmq, redisClient redis.Client, token func() string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestType := mux.Vars(r)["requestType"]
		token := token()

		// Create redis record of response.
		response := &data.Response{
			Body:        "{}",
			ID:          token,
			RequestType: requestType,
			Status:      "IN_PROGRESS",
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
			fmt.Sprintf("response:%s", token),
			responseJSON,
			0,
		)

		// TODO: enqueue message to perform the work

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
