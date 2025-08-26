package data

import "strings"

// Request type stuff
type RequestType int

const (
	NOOP RequestType = iota
	DEBUG
	DOWNLOAD
	// DEFINE NEW REQUEST TYPES HERE
)

var requestTypeMap = map[string]RequestType{
	"debug":    DEBUG,
	"noop":     NOOP,
	"download": DOWNLOAD,
}

func GetRequestType(rawRequestType string) RequestType {
	return requestTypeMap[strings.ToLower(rawRequestType)]
}

// Reqeusts to server
type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateRequest struct {
	Status string `json:"status"`
}
