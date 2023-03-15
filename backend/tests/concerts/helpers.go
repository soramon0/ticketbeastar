package concerts

import (
	"context"
	"testing"
	"ticketbeastar/pkg/models"
	"time"

	"github.com/uptrace/bun"
)

func SetupTable(t *testing.T, db *bun.DB) {
	_, err := db.NewCreateTable().Model((*models.Concert)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("NewCreateTable(Concert) err %v; want nil", err)
	}
}

func TeardownTable(t *testing.T, db *bun.DB) {
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
