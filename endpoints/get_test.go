package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/redis"
)

func setupGet(requestType data.RequestType, identifier string, status data.Status, body string) (*httptest.ResponseRecorder, *mux.Router) {
	recorder := httptest.NewRecorder()
	db, mock := redismock.NewClientMock()
	mockedRedis := &redis.Redis{
		Client: *db,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := mux.NewRouter()
	router.HandleFunc(
		"/get/{requestType}/{id}",
		Get(mockedRedis, logger),
	)

	// Set up fake data in mock redis.
	responseString, err := json.Marshal(&data.Record{
		ID:        identifier,
		Status:    status,
		Timestamp: 1607212800,
		Body:      body,
	})
	if err != nil {
		errors.New(
			fmt.Sprintf(
				"Error marshalling JSON: %s",
				err.Error(),
			),
		)
		return nil, nil
	}
	mock.ExpectGet(
		fmt.Sprintf(
			"%s:%s",
			requestType.String(),
			identifier,
		),
	).SetVal(
		string(responseString),
	)

	return recorder, router
}

func TestGetDone(t *testing.T) {
	requestType := data.NOOP
	recorder, router := setupGet(
		requestType,
		"doneID",
		data.DONE,
		"{\"response\":\"content\"}",
	)

	// Create request.
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"/get/%s/doneID",
			requestType.String(),
		),
		nil,
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
		"{\"body\":\"{\\\"response\\\":\\\"content\\\"}\",\"id\":\"doneID\",\"request_type\":\"NOOP\",\"status\":\"DONE\",\"timestamp\":1607212800}\n",
		recorder.Body.String(),
	)
}

func TestGetNotReady(t *testing.T) {
	requestType := data.NOOP
	recorder, router := setupGet(
		requestType,
		"in_progress",
		data.IN_PROGRESS,
		"{notdoneyet;jsongibberish}",
	)

	// Create request.
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"/get/%s/in_progress",
			requestType.String(),
		),
		nil,
	)
	require.NoError(t, err)

	// Run request.
	router.ServeHTTP(recorder, req)

	// Verify response.
	assert.Equal(
		t,
		recorder.Code,
		http.StatusOK,
	)
	assert.Equal(
		t,
		"{\"body\":\"{notdoneyet;jsongibberish}\",\"id\":\"in_progress\",\"request_type\":\"NOOP\",\"status\":\"IN_PROGRESS\",\"timestamp\":1607212800}\n",
		recorder.Body.String(),
	)
}
