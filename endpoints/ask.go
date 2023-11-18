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
		rawRequestType := mux.Vars(r)["requestType"]
		requestType := data.GetRequestType(rawRequestType)

		token := token()

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf(
				"error reading request body: %s\n",
				fmt.Errorf("%w", err),
			)
		}

		// Create redis record of request.
		record := &data.Record{
			Body:        string(requestBody),
			ID:          token,
			RequestType: requestType,
			Status:      data.IN_PROGRESS,
			Timestamp:   clock.Now().Unix(),
		}
		recordJSON, err := json.Marshal(record)
		if err != nil {
			log.Fatalf(
				"error marshalling JSON for: %v\nerr: %s\n",
				record,
				fmt.Errorf("%w", err),
			)
		}
		ctx := context.Background()

		// Put it into redis.
		key := fmt.Sprintf(
			"%s:%s",
			requestType.String(),
			token,
		)
		err = redisClient.Set(
			ctx,
			key,
			recordJSON,
		)
		if err != nil {
			log.Fatalf(
				"redis write failed for: %s\n%s\n%s\n",
				key,
				recordJSON,
				fmt.Errorf("%w", err),
			)
		}

		// Enqueue message to perform the work
		rabbitmq.Publish(string(recordJSON))

		// Debug message.
		log.Printf(
			"ask endpoint requested. [requestType=%s]",
			requestType.String(),
		)

		json.NewEncoder(w).Encode(record.ToAskResponse())
	}
}
