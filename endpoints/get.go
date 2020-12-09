package endpoints

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"

	"github.com/maxgarvey/cq_server/data"
)

// Get a response.
func Get(redisConnection redis.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse response id from URL.
		requestID := mux.Vars(r)["id"]
		log.Printf("requestID: %s", requestID)
		// Retrieve raw response from DB.
		rawResponse, err := redis.Values(
			redisConnection.Do(
				"HGETALL",
				fmt.Sprintf("response:%s", requestID)))
		if err != nil {
			log.Fatal(err)
		}

		// Parse response.
		var response data.Response
		err = redis.ScanStruct(rawResponse, &response)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(
			"get endpoint requested. [ID=%s, status=%s]",
			response.ID,
			response.Status)
		// If response is not ready.
		if response.Status != "DONE" {
			fmt.Fprintf(w, "not ready")
			return
		}

		// If response is ready, return it.
		fmt.Fprintf(w, response.Body)
	}
}
