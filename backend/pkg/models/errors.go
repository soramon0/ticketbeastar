package models

const (
	ErrInvalidPaymentToken = modelError("invalid payment token")
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}
