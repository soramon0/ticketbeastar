package users

import (
	"context"
	"testing"
	"ticketbeastar/pkg/models"

	"github.com/go-faker/faker/v4"
	"github.com/uptrace/bun"
)

func SetupTable(t *testing.T, db *bun.DB) {
	_, err := db.NewCreateTable().Model((*models.User)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("NewCreateTable(User) err %v; want nil", err)
	}
}

func TeardownTable(t *testing.T, db *bun.DB) {
	_, err := db.NewDropTable().Model((*models.User)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("Drop users table err %v; want nil", err)
	}
}

func CreateUsers(t *testing.T, db *bun.DB, size uint) *[]models.User {
	var users []models.User
	var i uint
	for i = 1; i <= size; i++ {
		user := CreateUser(t, db, false)
		users = append(users, *user)
	}

	_, err := db.NewInsert().Model(&users).Exec(context.Background())
	if err != nil {
		t.Fatalf("createUsers() err %v; want nil", err)
	}
	return &users
}

func CreateUser(t *testing.T, db *bun.DB, insert bool) *models.User {
	user := &models.User{Name: faker.Name()}
	if insert {
		_, err := db.NewInsert().Model(user).Exec(context.Background())
		if err != nil {
			t.Fatalf("createUser() err %v; want nil", err)
		}
	}
	return user
}
