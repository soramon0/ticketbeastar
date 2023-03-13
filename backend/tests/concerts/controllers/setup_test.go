package controllers_test

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"
	"ticketbeastar/pkg/configs"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/routes"
	"ticketbeastar/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type testServer struct {
	app *fiber.App
	db  *bun.DB
	cs  models.ConcertService
	log *log.Logger
}

func newTestServer() *testServer {
	logger := utils.InitLogger()
	app := fiber.New(configs.FiberConfig())
	db := database.OpenConnection(utils.GetTestDatabaseURL(), logger)
	services := models.NewServices(db)
	routes.Register(app, services, logger)

	return &testServer{
		app: app,
		db:  db,
		cs:  services.Concert,
		log: logger,
	}
}

func (ts *testServer) setup(t *testing.T) {
	_, err := ts.db.NewCreateTable().Model((*models.Concert)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("NewCreateTable(Concert) err %v; want nil", err)
	}
}

func (ts *testServer) teardown(t *testing.T) {
	defer database.CloseConnection(ts.db)

	_, err := ts.db.NewDropTable().Model((*models.Concert)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("Drop concerts table err %v; want nil", err)
	}
}

func (ts *testServer) hitGetEndpoint(t *testing.T, endpoint string, wantStatusCode int, wantCount int, wantError string) *models.APIResponse[models.Concert] {
	req := httptest.NewRequest("GET", endpoint, nil)
	resp, err := ts.app.Test(req)
	if err != nil {
		t.Fatalf("GET %s err %v", endpoint, err)
	}
	if resp.StatusCode != wantStatusCode {
		t.Fatalf("statusCode = %d; want %d", resp.StatusCode, wantStatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	var apiResponse models.APIResponse[models.Concert]
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		t.Fatalf("apiResponse unmarshal() err %v; want nil", err)
	}
	if apiResponse.Count != wantCount {
		log.Fatalf("apiResponse.Count mismatch. want %d; got %d", wantCount, apiResponse.Count)
	}
	if wantError == "" && apiResponse.Error != nil {
		log.Fatalf("apiResponse.Error got %v; want nil", apiResponse.Error)
	}
	if apiResponse.Error != nil {
		if apiResponse.Error.Message != wantError {
			log.Fatalf("apiResponse.Error got %q; want %q", apiResponse.Error.Message, wantError)
		}
	}
	return &apiResponse
}
