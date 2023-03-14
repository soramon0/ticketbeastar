package controllers_test

import (
	"context"
	"database/sql"
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
		defer ts.teardown(t)
		t.Fatalf("NewCreateTable(Concert) err %v; want nil", err)
	}
}

func (ts *testServer) teardown(t *testing.T) {
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
		t.Fatalf("apiResponse.Count mismatch. want %d; got %d", wantCount, apiResponse.Count)
	}
	if wantError == "" && apiResponse.Error != nil {
		t.Fatalf("apiResponse.Error got %v; want nil", apiResponse.Error)
	}
	if apiResponse.Error != nil {
		if apiResponse.Error.Message != wantError {
			t.Fatalf("apiResponse.Error got %q; want %q", apiResponse.Error.Message, wantError)
		}
	}
	return &apiResponse
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
