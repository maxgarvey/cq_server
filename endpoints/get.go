package endpoints

import (
	"fmt"
	"log"
	"net/http"
)

// Get a response.
func Get(w http.ResponseWriter, r *http.Request) {
	// TODO: IMPLEMENT LOGIC
	fmt.Fprintf(w, "get")
	log.Printf("get endpoint requested")
}
