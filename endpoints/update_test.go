package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-redis/redismock/v9"
	"github.com/gorilla/mux"
	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUpdate(
	t *testing.T,
	requestType data.RequestType,
	identifier string,
	initialStatus data.Status,
	finalStatus data.Status,
) (*httptest.ResponseRecorder, *mux.Router) {
	recorder := httptest.NewRecorder()
	db, mock := redismock.NewClientMock()
	mockedRedis := &redis.Redis{
		Client: *db,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// 12/06/2020 @ 12:00am (UTC)
	timestamp, _ := time.Parse(
		"2006-01-02T15:04:05-0700",
		"2020-11-06T00:00:00-0000",
	)
	clock := clock.NewMock()
	clock.Set(timestamp)

	router := mux.NewRouter()
	router.HandleFunc(
		"/update/{requestType}/{id}",
		Update(clock, mockedRedis, logger),
	)

	initialRecord := data.Record{
		ID:        identifier,
		Status:    initialStatus,
		Timestamp: 1607212800,
		Body:      "{}",
	}

	finalRecord := data.Record{
		ID:        identifier,
		Status:    finalStatus,
		Timestamp: 1607212800,
		Body:      "{}",
	}

	// Set up fake data in mock redis.
	initialRecordString, err := json.Marshal(&initialRecord)
	require.NoError(t, err)
	redisToken := fmt.Sprintf(
		"%s:%s",
		requestType.String(),
		identifier,
	)
	mock.ExpectGet(
		redisToken,
	).SetVal(
		string(initialRecordString),
	)
	finalRecordString, err := json.Marshal(&finalRecord)
	require.NoError(t, err)
	mock.ExpectSet(
		redisToken,
		finalRecordString,
		0,
	).SetVal("OK")

	return recorder, router
}

func TestUpdateToDone(t *testing.T) {
	requestType := data.NOOP
	recorder, router := setupUpdate(
		t,
		requestType,
		"updateID",
		data.IN_PROGRESS,
		data.DONE,
	)

	// Create request.
	requestJSON, err := json.Marshal(&data.UpdateRequest{
		Status: data.DONE.String(),
	})
	require.NoError(
		t,
		err,
	)
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"/update/%s/updateID",
			requestType.String(),
		),
		bytes.NewReader(requestJSON),
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
		"{\"body\":\"{}\",\"id\":\"updateID\",\"request_type\":\"NOOP\",\"status\":\"DONE\",\"timestamp\":1607212800}\n",
		recorder.Body.String(),
	)
}
