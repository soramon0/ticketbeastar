package models

type APIResponse struct {
	Data  interface{} `json:"data"`
	Count int         `json:"count,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}
