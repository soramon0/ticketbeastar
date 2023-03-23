package models

const (
	ErrInvalidPaymentToken = modelError("invalid payment token")
	ErrNotEnoughTickets    = modelError("tickets not enough to fullfil request")
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}
