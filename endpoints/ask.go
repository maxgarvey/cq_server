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
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
)

// Ask enqueues a request and creates an entry in redis to track it.
func Ask(
	clock clock.Clock, rabbitmq rabbitmq.Rabbit, redisClient *redis.Redis, token func() string, logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		PerformAsk(
			w, r, clock, rabbitmq, redisClient, token, logger,
		)
	}
}

// Make this a function so we can reuse it for AdminAsk
func PerformAsk(
	w http.ResponseWriter,
	r *http.Request,
	clock clock.Clock,
	rabbitmq rabbitmq.Rabbit,
	redisClient *redis.Redis,
	token func() string,
	logger *slog.Logger,
) {
	rawRequestType := mux.Vars(r)["requestType"]
	requestType := data.GetRequestType(rawRequestType)

	thisToken := token()

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

	// Create redis record of request.
	record := &data.Record{
		Body:        string(requestBody),
		ID:          thisToken,
		RequestType: requestType,
		Status:      data.IN_PROGRESS,
		Timestamp:   clock.Now().Unix(),
	}
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

	// Put it into redis.
	key := fmt.Sprintf(
		"%s:%s",
		requestType.String(),
		thisToken,
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

	// Enqueue message to perform the work
	rabbitmq.Publish(string(recordJSON))

	// Debug message.
	logger.Debug(
		fmt.Sprintf(
			"ask endpoint requested. [requestType=%s]",
			requestType.String(),
		),
	)

	json.NewEncoder(w).Encode(record.ToAskResponse())
}
