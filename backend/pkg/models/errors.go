package models

const (
	ErrInvalidPaymentToken = modelError("invalid payment token")
	ErrNotEnoughTickets    = modelError("tickets not enough to fullfil request")
	ErrInvalidId           = modelError("invalid resouce id")
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}
