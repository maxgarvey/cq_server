package endpoints

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/gorilla/mux"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
)

func fakeRandomToken() string {
	return "token"
}

func setupAsk(requestType string) (*httptest.ResponseRecorder, *mux.Router, *rabbitmq.FakeRabbitmq) {
	recorder := httptest.NewRecorder()
	db, mock := redismock.NewClientMock()
	mockedRedis := &redis.Redis{
		Client: *db,
	}
	router := mux.NewRouter()
	// 12/06/2020 @ 12:00am (UTC)
	timestamp, _ := time.Parse(
		"2006-01-02T15:04:05-0700",
		"2020-11-06T00:00:00-0000",
	)
	clock := clockwork.NewFakeClockAt(timestamp)
	fakeRabbitmq := rabbitmq.InitFake()

	router.HandleFunc(
		"/ask/{requestType}",
		Ask(
			clock,
			&fakeRabbitmq,
			mockedRedis,
			fakeRandomToken,
		),
	)

	// Set up fake data in mock redis.
	record := &data.Record{
		Body:        "{}",
		ID:          "token",
		RequestType: "doWork",
		Status:      data.IN_PROGRESS,
		Timestamp:   clock.Now().Unix(),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		log.Fatal(err)
	}
	mock.ExpectSet(
		"doWork:token",
		recordJSON,
		0,
	)

	return recorder, router, &fakeRabbitmq
}

func TestAsk(t *testing.T) {
	// Prelim setup.
	recorder, router, fakeRabbitmq := setupAsk(
		"doWork",
	)

	// Set request body for incoming request to ask endpoint.
	requestBody := bytes.NewReader(
		[]byte(
			"{\"work\":\"content\"}",
		),
	)

	// Create request.
	req, err := http.NewRequest(
		"POST",
		"/ask/doWork",
		requestBody,
	)
	require.NoError(
		t,
		err,
	)

	// Run request.
	router.ServeHTTP(
		recorder,
		req,
	)

	// Verify response.
	assert.Equal(
		t,
		recorder.Code,
		http.StatusOK,
	)
	assert.Equal(
		t,
		"{\"id\":\"token\"}\n",
		recorder.Body.String(),
	)

	// Verify publish to rabbit.
	assert.Equal(
		t,
		fakeRabbitmq.PublishedMessages,
		[]string{
			"{\"Body\":\"{\\\"work\\\":\\\"content\\\"}\",\"ID\":\"token\",\"RequestType\":\"doWork\",\"Status\":0,\"Timestamp\":1604620800}",
		},
	)
}
