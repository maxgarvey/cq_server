package endpoints

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

// Get a response.
func Get(redisConnection redis.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := mux.Vars(r)["id"]
		log.Printf("get endpoint requested. [requestId=%s]", requestId)

		// TODO: IMPLEMENT LOGIC
		fmt.Fprintf(w, "get")
	}
}
