package controllers_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/tests"
	"ticketbeastar/tests/concerts"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestConcertsController(t *testing.T) {
	ts := tests.NewTestServer()
	defer database.CloseConnection(ts.Db)

	tests := map[string]func(t *testing.T){
		"can view single published concert": func(t *testing.T) {
			concert := concerts.CreateConcert(t, ts.Db, nil, "", true)
			resp := ts.Visit(t, fmt.Sprintf("/api/v1/concerts/%d", concert.Id))
			api := unmarshalConcert(t, resp.Body)

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
			concert := concerts.CreateConcert(t, ts.Db, &models.Concert{PublishedAt: sql.NullTime{}}, "", true)
			endpoint := fmt.Sprintf("/api/v1/concerts/%d", concert.Id)
			resp := ts.Visit(t, endpoint)
			api := unmarshalConcert(t, resp.Body)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusNotFound)
			ts.AssertResponseCount(t, api.Count, 0)
			ts.AssertResponseError(t, api.Error, &models.APIError{Message: "Concert not found"})

			if api.Data != nil {
				t.Fatalf("response data should be nil; got %v", api.Data)
			}
		},
		"can view list of published concerts": func(t *testing.T) {
			concerts.CreateConcert(t, ts.Db, &models.Concert{PublishedAt: sql.NullTime{}}, "", true)
			concert2 := concerts.CreateConcert(t, ts.Db, &models.Concert{PublishedAt: sql.NullTime{Time: time.Now(), Valid: true}}, "", true)
			resp := ts.Visit(t, "/api/v1/concerts")
			api := unmarshalConcerts(t, resp.Body)

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

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			concerts.SetupTable(t, ts.Db)
			defer concerts.TeardownTable(t, ts.Db)
			test(t)
		})
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
