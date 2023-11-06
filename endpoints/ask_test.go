package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/gorilla/mux"
	"github.com/jonboulle/clockwork"
	"github.com/maxgarvey/cq_server/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fakeRandomToken() string {
	return "token"
}

func setupAsk(requestType string, body string) (*httptest.ResponseRecorder, *mux.Router) {
	recorder := httptest.NewRecorder()
	db, mock := redismock.NewClientMock()
	router := mux.NewRouter()
	// 12/06/2020 @ 12:00am (UTC)
	timestamp, _ := time.Parse(
		"2006-01-02T15:04:05-0700",
		"2020-11-06T00:00:00-0000",
	)
	clock := clockwork.NewFakeClockAt(timestamp)

	router.HandleFunc(
		"/ask/{requestType}",
		Ask(
			clock,
			*db, fakeRandomToken,
		),
	)

	// Set up fake data in mock redis.
	response := &data.Response{
		Body:        "{}",
		ID:          "token",
		RequestType: "doWork",
		Status:      "IN_PROGRESS",
		Timestamp:   clock.Now().Unix(),
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	mock.ExpectSet(
		"response:token",
		responseJSON,
		0,
	)

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

	// Verify redis set transaction.
	// assert.Equal(t, 1, client.(command))
}
