package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type FakeRabbitmq struct {
	PublishedMessages []string
	ConsumeChannel    <-chan amqp.Delivery
}

func InitFake() FakeRabbitmq {
	return FakeRabbitmq{}
}

func (f FakeRabbitmq) Close() {}

func (f *FakeRabbitmq) Consume() <-chan amqp.Delivery {
	return f.ConsumeChannel
}

func (f *FakeRabbitmq) Publish(message string) {
	f.PublishedMessages = append(f.PublishedMessages, message)
}
