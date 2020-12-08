package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsk(t *testing.T) {
	// Prelim setup.
	recorder := httptest.NewRecorder()

	// Create request.
	req, err := http.NewRequest("GET", "/ask", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Process request.
	http.HandlerFunc(Ask).ServeHTTP(recorder, req)

	// Verify response.
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, recorder.Body.String(), "ask")
}
