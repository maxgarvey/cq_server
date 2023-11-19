package endpoints

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Basic healthcheck endpoint.
func Health(logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("health endpoint requested")
		fmt.Fprintf(
			w,
			"healthy",
		)
	}
}
