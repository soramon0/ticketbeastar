package models

import (
	"github.com/uptrace/bun"
)

type Services struct {
	User UserService
}

func NewServices(db *bun.DB) *Services {
	us := NewUserService(db)

	return &Services{
		User: us,
	}
}
