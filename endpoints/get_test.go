package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/gorilla/mux"
	"github.com/maxgarvey/cq_server/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGet(identifier string, status string, body string) (*httptest.ResponseRecorder, *mux.Router) {
	recorder := httptest.NewRecorder()
	db, mock := redismock.NewClientMock()
	router := mux.NewRouter()
	router.HandleFunc("/get/{id}", Get(*db))

	// Set up fake data in mock redis.
	responseString, err := json.Marshal(&data.Response{
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
			"response:%s",
			identifier,
		),
	).SetVal(
		string(responseString),
	)

	return recorder, router
}

func TestGetDone(t *testing.T) {
	recorder, router := setupGet(
		"doneID",
		"DONE",
		"{\"response\":\"content\"}",
	)

	// Create request.
	req, err := http.NewRequest(
		"GET",
		"/get/doneID",
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
		"{\"response\":\"content\"}",
		recorder.Body.String(),
	)
}

func TestGetNotReady(t *testing.T) {
	recorder, router := setupGet(
		"in_progress",
		"IN_PROGRESS",
		"{notdoneyet;jsongibberish}",
	)

	// Create request.
	req, err := http.NewRequest(
		"GET",
		"/get/in_progress",
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
		"not ready",
		recorder.Body.String(),
	)
}
