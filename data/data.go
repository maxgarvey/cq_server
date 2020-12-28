package data

// Response object in redis cache.
type Response struct {
	Body        string `redis:"body"`
	ID          string `redis:"id"`
	RequestType string `redis:"requestType"`
	Status      string `redis:"status"`
	Timestamp   int64  `redis:"timestamp"`
}

// AskResponse is response body to ask endpoint.
type AskResponse struct {
	ID string `json:"id"`
}
