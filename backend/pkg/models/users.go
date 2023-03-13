package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id   uint64 `bun:"id,pk,autoincrement" json:"id"`
	Name string `bun:"name,notnull" json:"name"`
}

type UserService interface {
	// Methods for querying users
	// ByID(id string) (*User, error)
	Find() (*[]User, error)
	// ByEmail(email string) (*User, error)

	// Methods for altering users
	// Create() (*User, error)
	// Update(user *User) error
	// Delete(id string) error
}

type userService struct {
	db *bun.DB
}

func NewUserService(db *bun.DB) UserService {
	return &userService{
		db: db,
	}
}

func (us *userService) Find() (*[]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	users := []User{}
	err := us.db.NewSelect().Model(&users).Scan(ctx)
	return &users, err
}
