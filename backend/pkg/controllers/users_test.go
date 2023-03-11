package controllers_test

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"
	"ticketbeastar/pkg/configs"
	"ticketbeastar/pkg/controllers"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"

	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func TestUsersController(t *testing.T) {
	ts := newTestServer()

	t.Run("List", func(t *testing.T) {
		ts.setup(t)
		defer ts.teardown(t)
		testUsersListing(t, ts)
	})

}

func testUsersListing(t *testing.T, ts *testServer) {
	users, err := createUsers(ts.db, 5)
	if err != nil {
		t.Fatalf("createUsers() err %v; want nil", err)
	}

	userC := controllers.NewUsers(ts.us, ts.log)
	ts.app.Get("/api/v1/users", userC.GetUsers)
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	resp, err := ts.app.Test(req)
	if err != nil {
		t.Fatal("GET /api/v1/users", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("statusCode = %d; want %d", resp.StatusCode, fiber.StatusOK)
	}
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var apiResponse models.APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		t.Fatalf("apiResponse unmarshal() err %v; want nil", err)
	}
	if apiResponse.Count != len(*users) {
		log.Fatalf("apiResponse.Count incorrect. want %d; got %d", apiResponse.Count, len(*users))
	}
	if apiResponse.Error != nil {
		log.Fatalf("apiResponse.Error got %v; want nil", apiResponse.Error)
	}
}

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
