package worker

import (
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/maxgarvey/cq_server/rabbitmq"
)

func setupWorkerTest() *rabbitmq.FakeRabbitmq {
	fakeRabbitmq := rabbitmq.InitFake()

	return &fakeRabbitmq
}

func testWorker() {
	fakeRabbitmq := setupWorkerTest()
	message := amqp.Delivery{
		Body: []byte{},
	}

	fakeRabbitmq.ConsumeChannel <- message
}
