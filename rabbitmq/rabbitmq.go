package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbitmq struct {
	Username   string
	Password   string
	Host       string
	Port       int
	Channel    *amqp.Channel
	Connection *amqp.Connection
	Queue      *amqp.Queue
}

func Init(username string, password string, host string, port int, queuename string) *Rabbitmq {
	// Connect to instance
	connection, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%d/",
			username,
			password,
			host,
			port,
		),
	)
	if err != nil {
		panic(err)
	}

	// Connect to channel
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	// Connect to queue
	queue, err := channel.QueueDeclare(
		queuename, // name of queue
		false,     // is durable?
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		panic(err)
	}

	return &Rabbitmq{
		Username:   username,
		Password:   password,
		Host:       host,
		Port:       port,
		Channel:    channel,
		Connection: connection,
		Queue:      &queue,
	}
}

func (r Rabbitmq) Publish(message string) {
	err := r.Channel.PublishWithContext(
		context.TODO(),
		"",
		r.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		// TODO: improve handling of error during publish
		panic(err)
	}
}

func (r Rabbitmq) Consume() <-chan amqp.Delivery {
	messages, err := r.Channel.ConsumeWithContext(
		context.TODO(),
		r.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		// TODO: improve handling of error during consume
		panic(err)
	}
	return messages
}
