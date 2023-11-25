package data

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

// Admin stuff:
type AdminLoginResponse struct {
	Token string `json:"token"`
}
