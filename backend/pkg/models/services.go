package models

import (
	"github.com/uptrace/bun"
)

type Services struct {
	User    UserService
	Concert ConcertService
}

func NewServices(db *bun.DB) *Services {
	us := NewUserService(db)
	cs := NewConcertService(db)

	return &Services{
		User:    us,
		Concert: cs,
	}
}
