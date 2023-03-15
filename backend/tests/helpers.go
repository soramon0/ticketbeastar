package tests

import (
	"context"
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

	"github.com/go-faker/faker/v4"
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

func SetupConcertTable(t *testing.T, db *bun.DB) {
	_, err := db.NewCreateTable().Model((*models.Concert)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("NewCreateTable(Concert) err %v; want nil", err)
	}
}

func TeardownConcertTable(t *testing.T, db *bun.DB) {
	_, err := db.NewDropTable().Model((*models.Concert)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("Drop concerts table err %v; want nil", err)
	}
}

func CreateConcert(t *testing.T, db *bun.DB, overrides *models.Concert, dateStr string, insert bool) *models.Concert {
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
		PublishedAt:           bun.NullTime{Time: time.Now()},
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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := db.NewInsert().Model(&concert).Exec(ctx)
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
	if c.PublishedAt.Time.IsZero() {
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

func SetupUserTable(t *testing.T, db *bun.DB) {
	_, err := db.NewCreateTable().Model((*models.User)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("NewCreateTable(User) err %v; want nil", err)
	}
}

func TeardownUserTable(t *testing.T, db *bun.DB) {
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
