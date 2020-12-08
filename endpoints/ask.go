package endpoints

import (
	"fmt"
	"log"
	"net/http"
)

// Make a request.
func Ask(w http.ResponseWriter, r *http.Request) {
	// TODO: IMPLEMENT LOGIC
	fmt.Fprintf(w, "ask")
	log.Printf("ask endpoint requested")
}
