package endpoints

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maxgarvey/cq_server/admin"
	"github.com/maxgarvey/cq_server/data"
	"github.com/maxgarvey/cq_server/redis"
)

func AdminGet(
	admin admin.Adminer, redisClient redis.Redis, logger slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the requestType from URL.
		rawRequestType := mux.Vars(r)["requestType"]
		requestType := data.GetRequestType(rawRequestType)

		// Parse response id from URL.
		requestID := mux.Vars(r)["id"]

		session := r.Header.Get("SESSION")
		valid, err := admin.ValidateSession(session)
		if err != nil || !valid {
			logger.Error(
				fmt.Sprintf(
					"Unable to validate session token: %s\n",
					fmt.Errorf("%w", err),
				),
			)
			return
		}

		err = admin.ExtendSession(session)
		if err != nil {
			logger.Error(
				fmt.Sprintf(
					"Error extending sessino: %s\n",
					fmt.Errorf("%w", err),
				),
			)
			return
		}

		record, err := GetFromRedis(
			requestID, requestType, redisClient,
		)
		if err != nil {
			logger.Error(
				fmt.Sprintf(
					"error retrieving record from redis: %s\n",
					fmt.Errorf("%w", err),
				),
			)
			return
		}

		// If response is ready, return it.
		json.NewEncoder(w).Encode(record.ToGetResponse())
	}
}
