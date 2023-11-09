package worker

import (
	"fmt"

	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
)

type Worker struct {
	Rabbitmq *rabbitmq.Rabbitmq
	Redis    *redis.Redis
}

func Init(rabbitmq *rabbitmq.Rabbitmq, redis *redis.Redis) *Worker {
	return &Worker{
		Rabbitmq: rabbitmq,
		Redis:    redis,
	}
}

func (w Worker) Work() {
	go func() {
		for msg := range w.Rabbitmq.Consume() {
			fmt.Printf("Received Message: %s\n", msg.Body)
			// TODO: handle message: do work, update redis record.
		}
	}()

	fmt.Println("Waiting for messages...")
}
