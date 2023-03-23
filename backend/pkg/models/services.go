package models

import (
	"strconv"

	"github.com/uptrace/bun"
)

type Services struct {
	User    UserService
	Concert ConcertService
	Order   OrderService
	Ticket  TicketService
}

func NewServices(db *bun.DB) *Services {
	us := NewUserService(db)
	cs := NewConcertService(db)
	os := NewOrderService(db)
	ts := NewTicketService(db)

	return &Services{
		User:    us,
		Concert: cs,
		Order:   os,
		Ticket:  ts,
	}
}

func parseId(id string) (uint64, error) {
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, ErrInvalidId
	}
	return parsedId, nil
}
