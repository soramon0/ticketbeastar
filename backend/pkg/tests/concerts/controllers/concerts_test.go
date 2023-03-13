package controllers_test

import (
	"testing"
	"ticketbeastar/pkg/controllers"
	"ticketbeastar/pkg/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestConcertsController(t *testing.T) {
	ts := newTestServer()

	t.Run("List", func(t *testing.T) {
		ts.setup(t)
		defer ts.teardown(t)
		testConcertListing(t, ts)
	})
}

func testConcertListing(t *testing.T, ts *testServer) {
	date, err := time.Parse(time.RFC822, "02 Dec 06 08:00 MST")
	if err != nil {
		t.Fatalf("time.Parse() err %v; want nil", err)
	}
	concert := models.Concert{
		Title:                 "The Red Chord",
		Subtitle:              "with Animosity and Lethargy",
		Date:                  date,
		TicketPrice:           3250,
		Venue:                 "The Mosh Pit",
		VenueAddress:          "123 Example Lane",
		City:                  "Golang city",
		State:                 "On",
		Zip:                   "17916",
		AdditionalInformation: "For tickets, call (555) 555-5555",
	}
	err = ts.cs.Create(&concert)
	if err != nil {
		t.Fatalf("Create(concert) err = %v, want nil", err)
	}

	concertsC := controllers.NewConcerts(ts.cs, ts.log)
	endpoint := "/api/v1/concerts"
	ts.app.Get(endpoint, concertsC.GetConcerts)

	ts.hitGetEndpoint(t, endpoint, fiber.StatusOK, 1, "")
}
