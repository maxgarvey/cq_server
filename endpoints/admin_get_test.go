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
	"github.com/maxgarvey/cq_server/admin"
	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAdminGet(
	requestType data.RequestType, identifier string, status data.Status, body string,
) (*httptest.ResponseRecorder, *mux.Router, *admin.MockAdmin) {
	recorder := httptest.NewRecorder()
	admin := admin.InitMock()
	logger := *slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, mock := redismock.NewClientMock()
	mockedRedis := &redis.Redis{
		Client: *db,
	}

	router := mux.NewRouter()
	router.HandleFunc(
		"/admin/get/{requestType}/{id}",
		AdminGet(
			&admin,
			*mockedRedis,
			logger,
		),
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
		return nil, nil, nil
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

	return recorder, router, &admin
}

func TestAdminGet(t *testing.T) {
	requestType := data.NOOP
	identifier := "doneID"
	recorder, router, mockAdmin := setupAdminGet(
		requestType,
		identifier,
		data.DONE,
		"{\"response\":\"content\"}",
	)
	// Create request.
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"/admin/get/%s/%s",
			requestType.String(),
			identifier,
		),
		nil,
	)
	req.Header.Set("SESSION", "token")
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

	// Verify mock
	assert.Equal(
		t,
		[]admin.ValidateSessionCall{
			{
				Token: "token",
			},
		},
		mockAdmin.ValidateSessionCalls,
	)
	assert.Equal(
		t,
		[]admin.ExtendSessionCall{
			{
				Token: "token",
			},
		},
		mockAdmin.ExtendSessionCalls,
	)
}
