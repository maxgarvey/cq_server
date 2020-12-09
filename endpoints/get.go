package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

// Response object in redis cache.
type Response struct {
	Body      string `json:"body"`
	ID        string `json:"id"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// Get a response.
func Get(redisConnection redis.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse response id from URL.
		requestID := mux.Vars(r)["id"]

		// Retrieve raw response from DB.
		rawResponse, err := redis.StringMap(
			redisConnection.Do(
				"HGETALL",
				fmt.Sprintf("response:%s", requestID)))
		if err != nil {
			log.Fatal(err)
		}

		// Parse response.
		response, err := parseResponse(rawResponse)
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

func parseResponse(rawResponse map[string]string) (*Response, error) {
	var err error
	response := new(Response)
	response.Body = rawResponse["body"]
	response.ID = rawResponse["id"]
	response.Status = rawResponse["status"]
	response.Timestamp, err = strconv.ParseInt(rawResponse["timestamp"], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return response, nil
}
