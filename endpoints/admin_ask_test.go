package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-redis/redismock/v9"
	"github.com/gorilla/mux"
	"github.com/maxgarvey/cq_server/admin"
	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAdminAsk(
	requestType data.RequestType, identifier string, status data.Status, body string,
) (*httptest.ResponseRecorder, *mux.Router, *admin.MockAdmin, *rabbitmq.FakeRabbitmq) {
	recorder := httptest.NewRecorder()
	admin := admin.InitMock()
	logger := *slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, mock := redismock.NewClientMock()
	mockedRedis := &redis.Redis{
		Client: *db,
	}

	// 12/06/2020 @ 12:00am (UTC)
	timestamp, _ := time.Parse(
		"2006-01-02T15:04:05-0700",
		"2020-11-06T00:00:00-0000",
	)
	clock := clock.NewMock()
	clock.Set(timestamp)
	fakeRabbitmq := rabbitmq.InitFake()

	router := mux.NewRouter()
	router.HandleFunc(
		"/admin/ask/{requestType}",
		AdminAsk(
			&admin,
			clock,
			&fakeRabbitmq,
			mockedRedis,
			fakeRandomToken,
			&logger,
		),
	)

	// Set up fake data in mock redis.
	record := &data.Record{
		Body:        "{\"work\":\"content\"}",
		ID:          "token",
		RequestType: requestType,
		Status:      data.IN_PROGRESS,
		Timestamp:   clock.Now().Unix(),
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		log.Fatalf(
			"error marshalling JSON: %s\n",
			fmt.Errorf("%w", err),
		)
	}
	mock.ExpectSet(
		fmt.Sprintf(
			"%s:token",
			requestType.String(),
		),
		recordJSON,
		0,
	).SetVal("OK")

	return recorder, router, &admin, &fakeRabbitmq
}

func TestAdminAsk(t *testing.T) {
	// Prelim setup.
	recorder, router, mockAdmin, fakeRabbitmq := setupAdminAsk(
		data.NOOP,
		"ask_id",
		data.DONE,
		"",
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
		"/admin/ask/NOOP",
		requestBody,
	)
	require.NoError(
		t,
		err,
	)
	req.Header.Add("SESSION", "session_token")

	// Run request.
	router.ServeHTTP(
		recorder,
		req,
	)

	// Verify response.
	assert.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)
	assert.Equal(
		t,
		"{\"id\":\"token\"}\n",
		recorder.Body.String(),
	)

	// Verify session authentication
	assert.Equal(
		t,
		[]admin.ValidateSessionCall{
			{
				Token: "session_token",
			},
		},
		mockAdmin.ValidateSessionCalls,
	)
	assert.Equal(
		t,
		[]admin.ExtendSessionCall{
			{
				Token: "session_token",
			},
		},
		mockAdmin.ExtendSessionCalls,
	)

	// Verify publish to rabbit.
	assert.Equal(
		t,
		[]string{
			"{\"body\":\"{\\\"work\\\":\\\"content\\\"}\",\"id\":\"token\",\"request_type\":0,\"status\":0,\"timestamp\":1604620800}",
		},
		fakeRabbitmq.PublishedMessages,
	)
}
