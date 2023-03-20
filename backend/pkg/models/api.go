package models

type possibleReturns interface {
	*User | *[]User | *Concert | *[]Concert | *Order | *[]Order | any
}

type APIResponse[T possibleReturns] struct {
	Data  T         `json:"data"`
	Count int       `json:"count,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

func NewAPIResponse[T possibleReturns](data T, count int) APIResponse[T] {
	return APIResponse[T]{
		Data:  data,
		Count: count,
	}
}

type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

type APIFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIValidaitonErrors struct {
	Errors []APIFieldError `json:"errors"`
}
