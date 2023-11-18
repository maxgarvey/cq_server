package data

import (
	"encoding/json"
	"log"
)

// Progress tatuses for a particular message indicating level of doneness.
type Status int

const (
	IN_PROGRESS Status = iota
	DONE
)

type RequestType int

const (
	NOOP RequestType = iota
	DEBUG
	// DEFINE NEW REQUEST TYPES HERE
)

var requestTypeMap = map[string]RequestType{
	"debug": DEBUG,
	"noop":  NOOP,
}

func GetRequestType(rawRequestType string) RequestType {
	return requestTypeMap[rawRequestType]
}

// Record object in redis cache.
type Record struct {
	Body        string      `redis:"body" json:"body" yaml:"body"`
	ID          string      `redis:"id" json:"id" yaml:"id"`
	RequestType RequestType `redis:"requestType" json:"request_type" yaml:"request_type"`
	Status      Status      `redis:"status" json:"status" yaml:"status"`
	Timestamp   int64       `redis:"timestamp" json:"timestamp" yaml:"timestamp"`
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

func (r *Record) DoWork() error {
	switch r.RequestType {
	// Print record for DEBUG type.
	case DEBUG:
		log.Default().Printf("Received record: %v\n", r)
	// Perform no options for NOOP.
	case NOOP:
	// Catch all also do nothing. This shouldn't be called.
	default:
	}

	return nil
}

// AskResponse is the response body to ask endpoint.
type AskResponse struct {
	ID string `json:"id"`
}

// Helper to render ask response.
func (r *Record) ToAskResponse() AskResponse {
	return AskResponse{
		ID: r.ID,
	}
}

// GetResponse is the response body to the get endpoint.
type GetResponse struct {
	Body        string `json:"body" yaml:"body"`
	ID          string `json:"id" yaml:"id"`
	RequestType string `json:"request_type" yaml:"request_type"`
	Status      string `json:"status" yaml:"status"`
	Timestamp   int64  `json:"timestamp" yaml:"timestamp"`
}

// Helper to render get response
func (r *Record) ToGetResponse() GetResponse {
	return GetResponse{
		Body:        r.Body,
		ID:          r.ID,
		RequestType: r.RequestType.String(),
		Status:      r.Status.String(),
		Timestamp:   r.Timestamp,
	}
}
