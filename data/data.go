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
	Body        string `redis:"body" json:"body" yaml:"body"`
	ID          string `redis:"id" json:"id" yaml:"id"`
	RequestType string `redis:"requestType" json:"request_type" yaml:"request_type"`
	Status      Status `redis:"status" json:"status" yaml:"status"`
	Timestamp   int64  `redis:"timestamp" json:"timestamp" yaml:"timestamp"`
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
