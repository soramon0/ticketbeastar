package tests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"ticketbeastar/pkg/configs"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/routes"
	"ticketbeastar/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type TestServer struct {
	App     *fiber.App
	Db      *bun.DB
	Service *models.Services
	Log     *log.Logger
}

func NewTestServer() *TestServer {
	logger := utils.InitLogger()
	app := fiber.New(configs.FiberConfig())
	db := database.OpenConnection(utils.GetTestDatabaseURL(), logger)
	services := models.NewServices(db)
	routes.Register(app, services, logger)

	return &TestServer{
		App:     app,
		Db:      db,
		Service: services,
		Log:     logger,
	}
}

func (ts *TestServer) Visit(t *testing.T, endpoint string) *http.Response {
	req := httptest.NewRequest("GET", endpoint, nil)
	resp, err := ts.App.Test(req)
	if err != nil {
		t.Fatalf("could not reach %s; err %v", endpoint, err)
	}

	return resp
}

func (ts *TestServer) AssertResponseStatus(t *testing.T, gotCode, wantCode int) {
	if gotCode != wantCode {
		t.Fatalf("response status code should be %d; got %d", wantCode, gotCode)
	}
}

func (ts *TestServer) AssertResponseError(t *testing.T, gotError, wantError *models.APIError) {
	if !reflect.DeepEqual(gotError, wantError) {
		t.Fatalf("api response error should %v; got %v", wantError, gotError)
	}
}

func (ts *TestServer) AssertResponseCount(t *testing.T, gotCount, wantCount int) {
	if gotCount != wantCount {
		t.Fatalf("api response count should be %d; got %d", wantCount, gotCount)
	}
}
