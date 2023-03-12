package controllers_test

import (
	"context"
	"log"
	"testing"
	"ticketbeastar/pkg/configs"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"

	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type testServer struct {
	app *fiber.App
	db  *bun.DB
	us  models.UserService
	log *log.Logger
}

func newTestServer() *testServer {
	logger := utils.InitLogger()
	app := fiber.New(configs.FiberConfig())
	db := database.OpenConnection(utils.GetTestDatabaseURL(), logger)
	us := models.NewUserService(db)

	return &testServer{
		app: app,
		db:  db,
		us:  us,
		log: logger,
	}
}

func (ts *testServer) setup(t *testing.T) {
	_, err := ts.db.NewCreateTable().Model((*models.User)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("NewCreateTable() err %v; want nil", err)
	}
}

func (ts *testServer) teardown(t *testing.T) {
	defer database.CloseConnection(ts.db)

	_, err := ts.db.NewDropTable().Model((*models.User)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("Drop users table err %v; want nil", err)
	}
}

func createUsers(db *bun.DB, size uint) (*[]models.User, error) {
	var users []models.User
	var i uint
	for i = 1; i <= size; i++ {
		user, err := createUser(db, false)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	_, err := db.NewInsert().Model(&users).Exec(context.Background())
	return &users, err
}

func createUser(db *bun.DB, insert bool) (*models.User, error) {
	user := &models.User{Name: faker.Name()}

	if insert {
		res, err := db.NewInsert().Model(user).Exec(context.Background())
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		user.Id = id
	}

	return user, nil
}
