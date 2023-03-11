package controllers_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"ticketbeastar/pkg/configs"
	"ticketbeastar/pkg/controllers"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
)

func TestUsersController(t *testing.T) {
	ts := newTestServer()

	t.Run("List", func(t *testing.T) {
		ts.setup(t)
		err := loadUsersFixtures(ts.db)
		if err != nil {
			t.Fatalf("loadUsersFixtures() err %v; want nil", err)
		}

		defer ts.teardown(t)

		testUsersListing(t, ts)
	})

}

func testUsersListing(t *testing.T, ts *testServer) {
	userC := controllers.NewUsers(ts.us, ts.log)
	ts.app.Get("/api/v1/users", userC.GetUsers)
	// http.Request
	req := httptest.NewRequest("GET", "/api/v1/users", nil)

	resp, err := ts.app.Test(req)
	if err != nil {
		t.Fatal("GET /api/v1/users", err)
	}

	if resp.StatusCode == fiber.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body)) // => Hello, World!
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

func loadUsersFixtures(db *bun.DB) error {
	fixture := dbfixture.New(db)
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	dir := filepath.Join(cwd, "..", "migrations", "fixtures")
	if err := fixture.Load(context.Background(), os.DirFS(dir), "users.yaml"); err != nil {
		return err
	}
	return nil
}
