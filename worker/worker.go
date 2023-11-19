package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
)

// Worker is responsible for consuming messages from the rabbit queue
// and doing work based off of the messages.
type Worker struct {
	Rabbitmq rabbitmq.Rabbit
	Redis    *redis.Redis
	Logger   *slog.Logger
}

// Creates a new instance of the worker struct
func Init(rabbitmq rabbitmq.Rabbit, redis *redis.Redis, logger *slog.Logger) *Worker {
	return &Worker{
		Rabbitmq: rabbitmq,
		Redis:    redis,
		Logger:   logger,
	}
}

// Consume from rabbit, do work and update redis
func (w Worker) Work() {
	go func() {
		for msg := range w.Rabbitmq.Consume() {
			w.HandleMessage(
				msg,
			)
		}
	}()

	w.Logger.Info(
		"Waiting for messages...",
	)
}

// HandleMessage handles a single message, structured this way for
// easy of unit testing.
func (w Worker) HandleMessage(msg amqp.Delivery) {
	// Read message.
	w.Logger.Debug(
		fmt.Sprintf(
			"Received Message: %s\n",
			msg.Body,
		),
	)
	thisBody := data.Record{}
	json.Unmarshal(msg.Body, &thisBody)
	w.Logger.Debug(
		fmt.Sprintf(
			"Unmarshalled body: %v\n",
			thisBody,
		),
	)

	// Find redis record for this message in redis.
	token := fmt.Sprintf(
		"%s:%s",
		thisBody.RequestType.String(),
		thisBody.ID,
	)
	context := context.TODO()
	redisRecord, err := w.Redis.Get(
		context,
		token,
	)
	if err != nil {
		w.Logger.Error(
			fmt.Sprintf(
				"error requesting redis record. err=%v",
				err,
			),
		)
	}

	// TODO: do the work to handle the message.

	// Update the record to update Redis.
	redisRecord.Status = data.DONE

	// Create JSON to set in Redis
	updatedRedisRecordJSON, err := json.Marshal(redisRecord)
	if err != nil {
		w.Logger.Error(
			fmt.Sprintf(
				"error marshalling updated redis record. err=%v",
				err,
			),
		)
	}

	// Update Redis
	err = w.Redis.Set(
		context,
		token,
		updatedRedisRecordJSON,
	)
	if err != nil {
		w.Logger.Error(
			fmt.Sprintf(
				"error updating redis record. err=%v",
				err,
			),
		)
	}
}
