package models

import (
	"github.com/uptrace/bun"
)

type Services struct {
	User    UserService
	Concert ConcertService
	Order   OrderService
}

func NewServices(db *bun.DB) *Services {
	us := NewUserService(db)
	cs := NewConcertService(db)
	os := NewOrderService(db)

	return &Services{
		User:    us,
		Concert: cs,
		Order:   os,
	}
}
