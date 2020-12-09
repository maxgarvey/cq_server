package endpoints

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(identifier string, status string, body string) (*httptest.ResponseRecorder, *mux.Router) {
	recorder := httptest.NewRecorder()
	redisConnection := redigomock.NewConn()
	router := mux.NewRouter()
	router.HandleFunc("/get/{id}", Get(redisConnection))

	// Set up fake data in mock redis.
	redisConnection.Command(
		"HGETALL",
		fmt.Sprintf("response:%s", identifier)).ExpectMap(
		map[string]string{
			"id":        identifier,
			"status":    status,
			"timestamp": "1607212800", // 12/06/2020 @ 12:00am (UTC)
			"body":      body,
		})

	return recorder, router
}

func TestGetDone(t *testing.T) {
	recorder, router := setup("done", "DONE", "{\"response\":\"content\"}")

	// Create request.
	req, err := http.NewRequest("GET", "/get/done", nil)
	require.NoError(t, err)

	// Run request.
	router.ServeHTTP(recorder, req)

	// Verify response.
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(
		t, "{\"response\":\"content\"}", recorder.Body.String())
}

func TestGetNotReady(t *testing.T) {
	recorder, router := setup(
		"in_progress", "IN_PROGRESS", "{notdoneyet;jsongibberish}")

	// Create request.
	req, err := http.NewRequest("GET", "/get/in_progress", nil)
	require.NoError(t, err)

	// Run request.
	router.ServeHTTP(recorder, req)

	// Verify response.
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, recorder.Body.String(), "not ready")
}
