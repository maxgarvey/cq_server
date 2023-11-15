package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
)

type Worker struct {
	Rabbitmq rabbitmq.Rabbit
	Redis    *redis.Redis
}

func Init(rabbitmq rabbitmq.Rabbit, redis *redis.Redis) *Worker {
	return &Worker{
		Rabbitmq: rabbitmq,
		Redis:    redis,
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

	fmt.Println(
		"Waiting for messages...",
	)
}

func (w Worker) HandleMessage(msg amqp.Delivery) {
	// Read message.
	fmt.Printf("Received Message: %s\n", msg.Body)
	thisBody := data.Record{}
	json.Unmarshal(msg.Body, &thisBody)
	fmt.Printf("Unmarshalled body: %v\n", thisBody)

	// Find redis record for this message in redis.
	token := fmt.Sprintf("response:%s", thisBody.ID)
	context := context.TODO()
	redisRecord, err := w.Redis.Get(
		context,
		token,
	)
	if err != nil {
		log.Fatalf(
			"error requesting redis record. err=%v",
			err,
		)
	}

	// TODO: do the work to handle the message.

	// Update the record to update Redis.
	redisRecord.Status = data.DONE

	// Create JSON to set in Redis
	updatedRedisRecordJSON, err := json.Marshal(redisRecord)
	if err != nil {
		log.Fatalf(
			"error marshalling updated redis record. err=%v",
			err,
		)
	}

	// Update Redis
	err = w.Redis.Set(
		context,
		token,
		updatedRedisRecordJSON,
	)
	if err != nil {
		log.Fatalf(
			"error updating redis record. err=%v",
			err,
		)
	}
}
