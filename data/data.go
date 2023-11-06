package data

import (
	"encoding/json"
)

// Response object in redis cache.
type Response struct {
	Body        string `redis:"body"`
	ID          string `redis:"id"`
	RequestType string `redis:"requestType"`
	Status      string `redis:"status"`
	Timestamp   int64  `redis:"timestamp"`
}

func (r Response) MarshalBinary() ([]byte, error) {
	return json.Marshal(
		r,
	)
}
func (r *Response) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(
		data,
		r,
	)
}

// AskResponse is response body to ask endpoint.
type AskResponse struct {
	ID string `json:"id"`
}
