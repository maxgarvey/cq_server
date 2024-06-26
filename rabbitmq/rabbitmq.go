package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/maxgarvey/cq_server/config"
)

type Rabbit interface {
	Close()
	Consume() <-chan amqp.Delivery
	Publish(message string)
}

type Rabbitmq struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
	Queue      *amqp.Queue
}

func Init(config config.Rabbitmq) (*Rabbitmq, error) {
	// Connect to instance
	connection, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%d/",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
		),
	)
	if err != nil {
		return nil, err
	}

	// Connect to channel
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	// Connect to queue
	queue, err := channel.QueueDeclare(
		config.Queuename, // name of queue
		false,            // is durable?
		false,            // auto delete
		false,            // exclusive
		false,            // no wait
		nil,              // args
	)
	if err != nil {
		return nil, err
	}

	return &Rabbitmq{
		Channel:    channel,
		Connection: connection,
		Queue:      &queue,
	}, nil
}

// Close connection to RabbitMQ
func (r Rabbitmq) Close() {
	// TODO: is this the right order of operations?
	r.Channel.Close()
	r.Connection.Close()
}

// Enqueue a message to Rabbit MQ.
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

// Consume from the queue.
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
