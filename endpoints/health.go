package endpoints

import (
	"fmt"
	"log"
	"net/http"
)

// Basic healthcheck endpoint.
func Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "healthy")
	log.Printf("health endpoint requested")
}
