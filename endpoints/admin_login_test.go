package endpoints

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/maxgarvey/cq_server/admin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAdminLogin() (*httptest.ResponseRecorder, *mux.Router, *admin.MockAdmin) {
	recorder := httptest.NewRecorder()
	admin := admin.InitMock()
	logger := *slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := mux.NewRouter()
	router.HandleFunc(
		"/admin/login",
		AdminLogin(
			&admin,
			logger,
		),
	)

	return recorder, router, &admin
}

func TestAdminLogin(t *testing.T) {
	recorder, router, mockAdmin := setupAdminLogin()

	// Create request body
	requestBody := bytes.NewReader(
		[]byte(
			"{\"username\":\"my_user\",\"password\":\"my_password\"}",
		),
	)

	// Create request.
	req, err := http.NewRequest(
		"POST",
		"/admin/login",
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
		"{\"token\":\"session\"}\n",
		recorder.Body.String(),
	)

	// Verify mock call
	assert.Equal(
		t,
		[]admin.LoginCall{
			{
				Username: "my_user",
				Password: "my_password",
			},
		},
		mockAdmin.LoginCalls,
	)
}
