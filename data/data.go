package data

import (
	"encoding/json"
	"log"
	"time"

	"github.com/thanhpk/randstr"
)

// Progress tatuses for a particular message indicating level of doneness.
type Status int

const (
	IN_PROGRESS Status = iota
	DONE
)

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

// Method to make tokens for incoming requests.
func MakeToken() string {
	return randstr.String(20)
}

type User struct {
	ID        int
	Username  string
	CreatedAt string
	LastLogin string
}

type Session struct {
	UserID    int
	Token     string
	CreatedAt time.Time
	GoodUntil time.Time
}
