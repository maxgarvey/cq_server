package endpoints

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fakeRandomToken() string {
	return "token"
}

func setupAsk(requestType string, body string) (*httptest.ResponseRecorder, *mux.Router) {
	recorder := httptest.NewRecorder()
	redisConnection := redigomock.NewConn()
	router := mux.NewRouter()
	router.HandleFunc(
		"/ask/{requestType}",
		Ask(redisConnection, fakeRandomToken))

	// Set up fake data in mock redis.
	redisConnection.Command(
		"SET",
		"token").ExpectMap(
		map[string]string{
			"id":        "token",
			"status":    "IN_PROGRESS",
			"timestamp": "1607212800", // 12/06/2020 @ 12:00am (UTC)
			"body":      body,
		})

	return recorder, router
}

func TestAsk(t *testing.T) {
	// Prelim setup.
	recorder, router := setupAsk("doWork", "{\"work\":\"content\"}")

	// Create request.
	req, err := http.NewRequest("POST", "/ask/doWork", nil)
	require.NoError(t, err)

	// Run request.
	router.ServeHTTP(recorder, req)

	// Verify response.
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(
		t, "{\"id\":\"token\"}\n", recorder.Body.String())
}
