package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	// Prelim setup.
	recorder := httptest.NewRecorder()

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
	http.HandlerFunc(Health).ServeHTTP(
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
