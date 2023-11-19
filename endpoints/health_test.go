package endpoints

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	// Prelim setup.
	recorder := httptest.NewRecorder()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create request.
	req, err := http.NewRequest(
		"GET",
		"/health",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Process request.
	http.HandlerFunc(
		Health(logger),
	).ServeHTTP(
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
		recorder.Body.String(),
		"healthy",
	)
}
