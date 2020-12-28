package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jonboulle/clockwork"

	"github.com/maxgarvey/cq_server/data"
)

// Ask enqueues a request and creates an entry in redis to track it.
func Ask(clock clockwork.Clock, redisConnection redis.Conn, token func() string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestType := mux.Vars(r)["requestType"]
		token := token()

		// Create redis record of response.
		response := &data.Response{
			Body:        "{}",
			ID:          token,
			RequestType: requestType,
			Status:      "IN_PROGRESS",
			Timestamp:   clock.Now().Unix(),
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		// Put it into redis.
		redisConnection.Do(
			"SET",
			fmt.Sprintf("response:%s", token),
			responseJSON)

		// TODO: enqueue message to perform the work

		log.Printf("ask endpoint requested. [requestType=%s]", requestType)

		// Return token associated with this request.
		askResp := &data.AskResponse{
			ID: token,
		}
		json.NewEncoder(w).Encode(&askResp)
	}
}
