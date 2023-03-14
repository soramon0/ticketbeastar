package controllers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
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
	"time"

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
	_, err := ts.db.NewDropTable().Model((*models.Concert)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("Drop concerts table err %v; want nil", err)
	}
}

func (ts *testServer) visit(t *testing.T, endpoint string) *http.Response {
	req := httptest.NewRequest("GET", endpoint, nil)
	resp, err := ts.app.Test(req)
	if err != nil {
		t.Fatalf("could not reach %s; err %v", endpoint, err)
	}

	return resp
}

func (ts *testServer) createConcert(t *testing.T, overrides *models.Concert, dateStr string, insert bool) *models.Concert {
	if dateStr == "" {
		dateStr = "02 Dec 06 08:00 MST"
	}
	date, err := time.Parse(time.RFC822, dateStr)
	if err != nil {
		t.Fatalf("createConcert(date) err %v; want nil", err)
	}

	concert := models.Concert{
		Title:                 "The Red Chord",
		Subtitle:              "with Animosity and Lethargy",
		Date:                  date,
		PublishedAt:           sql.NullTime{Time: time.Now(), Valid: true},
		TicketPrice:           3250,
		Venue:                 "The Mosh Pit",
		VenueAddress:          "123 Example Lane",
		City:                  "Golang city",
		State:                 "On",
		Zip:                   "17916",
		AdditionalInformation: "For tickets, call (555) 555-5555",
	}
	if overrides != nil {
		overrideConcert(&concert, *overrides)
	}

	if insert {
		err = ts.cs.Create(&concert)
		if err != nil {
			t.Fatalf("Create(concert) err = %v, want nil", err)
		}
	}

	return &concert
}

func overrideConcert(concert *models.Concert, c models.Concert) {
	if c.Title != "" {
		concert.Title = c.Title
	}
	if c.Subtitle != "" {
		concert.Subtitle = c.Subtitle
	}
	if c.TicketPrice != 0 {
		concert.TicketPrice = c.TicketPrice
	}
	if !c.PublishedAt.Valid {
		concert.PublishedAt = c.PublishedAt
	}
	if c.Venue != "" {
		concert.Title = c.Venue
	}
	if c.VenueAddress != "" {
		concert.Title = c.VenueAddress
	}
	if c.City != "" {
		concert.Title = c.City
	}
	if c.State != "" {
		concert.Title = c.State
	}
	if c.Zip != "" {
		concert.Title = c.Zip
	}
	if c.AdditionalInformation != "" {
		concert.Title = c.AdditionalInformation
	}
}

func unmarshalConcert(t *testing.T, body io.ReadCloser) models.APIResponse[*models.Concert] {
	content, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("could not read response body; err %v", err)
	}
	defer body.Close()

	var resp models.APIResponse[*models.Concert]
	if err := json.Unmarshal(content, &resp); err != nil {
		t.Fatalf("could not unmarshal concert response body; err %v", err)
	}
	return resp
}

func unmarshalConcerts(t *testing.T, body io.ReadCloser) models.APIResponse[*[]models.Concert] {
	content, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("could not read response body; err %v", err)
	}
	defer body.Close()

	var resp models.APIResponse[*[]models.Concert]
	if err := json.Unmarshal(content, &resp); err != nil {
		t.Fatalf("could not unmarshal concerts response body; err %v", err)
	}
	return resp
}

func assertResponseStatus(t *testing.T, gotCode, wantCode int) {
	if gotCode != wantCode {
		t.Fatalf("response status code should be %d; got %d", wantCode, gotCode)
	}
}

func assertResponseError(t *testing.T, gotError, wantError *models.APIError) {
	if !reflect.DeepEqual(gotError, wantError) {
		t.Fatalf("api response error should %v; got %v", wantError, gotError)
	}
}

func assertResponseCount(t *testing.T, gotCount, wantCount int) {
	if gotCount != wantCount {
		t.Fatalf("api response count should be %d; got %d", wantCount, gotCount)
	}
}
