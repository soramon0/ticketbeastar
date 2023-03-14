package controllers_test

import (
	"database/sql"
	"fmt"
	"testing"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestConcertsController(t *testing.T) {
	ts := newTestServer()
	defer database.CloseConnection(ts.db)

	tests := map[string]func(t *testing.T, ts *testServer){
		"can view single published concert": func(t *testing.T, ts *testServer) {
			concert := ts.createConcert(t, nil, "", true)
			endpoint := fmt.Sprintf("/api/v1/concerts/%d", concert.Id)
			resp := unmarshalConcert(t, ts.visit(t, endpoint, fiber.StatusOK))
			assertResponse(t, resp, 0, "")

			if resp.Data == nil {
				t.Fatal("api response data should not be empty")
			}
			if resp.Data.Id != concert.Id {
				t.Fatalf("concert id mismatch want %d, got %d", concert.Id, resp.Data.Id)
			}
		},
		"cannot view single unpublished concert": func(t *testing.T, ts *testServer) {
			concert := ts.createConcert(t, &models.Concert{PublishedAt: sql.NullTime{}}, "", true)
			endpoint := fmt.Sprintf("/api/v1/concerts/%d", concert.Id)
			resp := unmarshalConcert(t, ts.visit(t, endpoint, fiber.StatusNotFound))
			assertResponse(t, resp, 0, "Concert not found")

			if resp.Data != nil {
				t.Fatalf("response data should be nil; got %v", resp.Data)
			}
		},
		"can view list of published concerts": func(t *testing.T, ts *testServer) {
			ts.createConcert(t, &models.Concert{PublishedAt: sql.NullTime{}}, "", true)
			concert2 := ts.createConcert(t, &models.Concert{PublishedAt: sql.NullTime{Time: time.Now(), Valid: true}}, "", true)
			resp := unmarshalConcerts(t, ts.visit(t, "/api/v1/concerts", fiber.StatusOK))
			assertResponse(t, resp, 1, "")

			data := *resp.Data
			if len(data) != 1 {
				t.Fatalf("response data should have one concert; got %d", len(data))
			}
			if data[0].Id != concert2.Id {
				t.Fatalf("concert id does not match; want %d, got %d", concert2.Id, data[0].Id)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ts.setup(t)
			defer ts.teardown(t)
			test(t, ts)
		})
	}
}
