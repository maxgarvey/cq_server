package data

import (
	"encoding/json"
)

type Status int

const (
	IN_PROGRESS Status = iota
	DONE
)

// Record object in redis cache.
type Record struct {
	Body        string `redis:"body"`
	ID          string `redis:"id"`
	RequestType string `redis:"requestType"`
	Status      Status `redis:"status"`
	Timestamp   int64  `redis:"timestamp"`
}

func (r Record) MarshalBinary() ([]byte, error) {
	return json.Marshal(
		r,
	)
}
func (r *Record) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(
		data,
		r,
	)
}

// AskResponse is response body to ask endpoint.
type AskResponse struct {
	ID string `json:"id"`
}
