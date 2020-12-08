package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	// Prelim setup.
	recorder := httptest.NewRecorder()

	// Create request.
	req, err := http.NewRequest("GET", "/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Run request.
	http.HandlerFunc(Get).ServeHTTP(recorder, req)

	// Verify response.
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, recorder.Body.String(), "get")
}
