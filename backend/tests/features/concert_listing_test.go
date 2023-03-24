package features_test

import (
	"fmt"
	"testing"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/tests"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func TestConcertsListing(t *testing.T) {
	ts := tests.NewTestServer(t)
	defer database.CloseConnection(ts.Db)

	testsCases := map[string]func(t *testing.T){
		"can view single published concert": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, nil, true)
			resp := ts.Visit(t, fmt.Sprintf("/api/v1/concerts/%d", concert.Id))
			api := tests.UnmarshalConcert(t, resp.Body)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusOK)
			ts.AssertResponseCount(t, api.Count, 0)
			ts.AssertResponseError(t, api.Error, nil)

			if api.Data == nil {
				t.Fatal("api response data should not be empty")
			}
			if api.Data.Id != concert.Id {
				t.Fatalf("concert id mismatch want %d, got %d", concert.Id, api.Data.Id)
			}
		},
		"cannot view single unpublished concert": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, &models.Concert{PublishedAt: bun.NullTime{}}, true)
			endpoint := fmt.Sprintf("/api/v1/concerts/%d", concert.Id)
			resp := ts.Visit(t, endpoint)
			api := tests.UnmarshalConcert(t, resp.Body)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusNotFound)
			ts.AssertResponseCount(t, api.Count, 0)
			ts.AssertResponseError(t, api.Error, &models.APIError{Message: "Concert not found"})

			if api.Data != nil {
				t.Fatalf("response data should be nil; got %v", api.Data)
			}
		},
		"can view list of published concerts": func(t *testing.T) {
			tests.CreateConcert(t, ts.Db, &models.Concert{PublishedAt: bun.NullTime{}}, true)
			concert2 := tests.CreateConcert(t, ts.Db, &models.Concert{PublishedAt: bun.NullTime{Time: time.Now()}}, true)
			resp := ts.Visit(t, "/api/v1/concerts")
			api := tests.UnmarshalConcerts(t, resp.Body)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusOK)
			ts.AssertResponseCount(t, api.Count, 1)
			ts.AssertResponseError(t, api.Error, nil)

			data := *api.Data
			if len(data) != 1 {
				t.Fatalf("response data should have one concert; got %d", len(data))
			}
			if data[0].Id != concert2.Id {
				t.Fatalf("concert id does not match; want %d, got %d", concert2.Id, data[0].Id)
			}
		},
	}

	for name, tc := range testsCases {
		t.Run(name, func(t *testing.T) {
			defer tests.TeardownConcertTable(t, ts.Db)
			tests.SetupConcertTable(t, ts.Db)
			tc(t)
		})
	}
}
