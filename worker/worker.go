package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
			fmt.Printf("Received Message: %s\n", msg.Body)
			thisBody := data.Record{}
			json.Unmarshal(msg.Body, &thisBody)
			fmt.Printf("Unmarshalled body: %v\n", thisBody)

			// TODO: handle message: do work

			// Update redis record.
			token := fmt.Sprintf("response:%s", thisBody.ID)
			context := context.TODO()
			redisRecord, err := w.Redis.Get(
				context,
				token,
			)
			// TODO: handle error in meaningful way.
			if err != nil {
				log.Fatalf(
					"error requesting redis record. err=%v",
					err,
				)
			}

			redisRecord.Status = data.DONE
			// Marshal updated record
			updatedRedisRecordJSON, err := json.Marshal(redisRecord)
			// TODO: handle error in meaningful way.
			if err != nil {
				log.Fatalf(
					"error marshalling updated redis record. err=%v",
					err,
				)
			}
			// Update redis record
			err = w.Redis.Set(
				context,
				token,
				updatedRedisRecordJSON,
			)
			// TODO: properly handle error.
			if err != nil {
				log.Fatalf(
					"error updating redis record. err=%v",
					err,
				)
			}
		}
	}()

	fmt.Println("Waiting for messages...")
}
