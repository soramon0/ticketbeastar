package models

type possibleReturns interface {
	*User | *[]User | *Concert | *[]Concert | any
}

type APIResponse[T possibleReturns] struct {
	Data  T         `json:"data"`
	Count int       `json:"count,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

func NewAPIResponse[T possibleReturns](data T, count int, err *APIError) APIResponse[T] {
	return APIResponse[T]{
		Data:  data,
		Count: count,
		Error: err,
	}
}

type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}
