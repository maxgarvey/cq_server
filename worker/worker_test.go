package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/go-redis/redismock/v9"
	"github.com/jonboulle/clockwork"

	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
)

func setupTestRabbit() (*Worker, redismock.ClientMock) {
	fakeRabbitmq := rabbitmq.InitFake()
	db, mock := redismock.NewClientMock()
	mockedRedis := &redis.Redis{
		Client: *db,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	worker := Init(
		&fakeRabbitmq,
		mockedRedis,
		logger,
	)

	return worker, mock
}

func TestWorker(t *testing.T) {
	timestamp, _ := time.Parse(
		"2006-01-02T15:04:05-0700",
		"2020-11-06T00:00:00-0000",
	)
	clock := clockwork.NewFakeClockAt(timestamp)

	requestType := data.NOOP

	worker, mockedRedis := setupTestRabbit()
	initialRawMessage := &data.Record{
		Body:        "{}",
		ID:          "token",
		RequestType: requestType,
		Status:      data.IN_PROGRESS,
		Timestamp:   clock.Now().Unix(),
	}
	initialMessageJson, err := json.Marshal(
		initialRawMessage,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Expect the get before work is done.
	mockedRedis.ExpectGet(
		fmt.Sprintf(
			"%s:token",
			requestType.String(),
		),
	).SetVal(
		string(
			initialMessageJson,
		),
	)
	finalRawMessage := initialRawMessage
	finalRawMessage.Status = data.DONE
	finalMessageJson, err := json.Marshal(
		finalRawMessage,
	)
	if err != nil {
		log.Fatal(err)
	}
	// Expect the set for when the updated record is written.
	mockedRedis.ExpectSet(
		fmt.Sprintf(
			"%s:token",
			requestType.String(),
		),
		finalMessageJson,
		0,
		// This is a response code from Redis, example here:
		// https://github.com/go-redis/redismock/blob/master/example/example.go#L78
	).SetVal("OK")

	message := amqp.Delivery{
		Body: initialMessageJson,
	}
	worker.HandleMessage(message)
}
