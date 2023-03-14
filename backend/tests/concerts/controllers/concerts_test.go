package controllers_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestConcertsController(t *testing.T) {
	ts := newTestServer()
	defer database.CloseConnection(ts.db)

	t.Run("List single concert", func(t *testing.T) {
		ts.setup(t)
		defer ts.teardown(t)
		testConcertListing(t, ts)
	})

	t.Run("List single public concert", func(t *testing.T) {
		ts.setup(t)
		defer ts.teardown(t)
		testPublicConcertListing(t, ts)
	})
}

func testConcertListing(t *testing.T, ts *testServer) {
	// Arrange
	concert := ts.createConcert(t, nil, "", true)
	endpoint := fmt.Sprintf("/api/v1/concerts/%d", concert.Id)

	// Act/Assert
	r := ts.hitGetEndpoint(t, endpoint, fiber.StatusOK, 0, "")

	// Assert
	if r.Data.Id != concert.Id {
		t.Fatalf("concert.id mismatch want %d, got %d", concert.Id, r.Data.Id)
	}
}

func testPublicConcertListing(t *testing.T, ts *testServer) {
	ts.createConcert(t, &models.Concert{PublishedAt: sql.NullTime{}}, "", true)
	concert2 := ts.createConcert(t, &models.Concert{PublishedAt: sql.NullTime{Time: time.Now(), Valid: true}}, "", true)

	req := httptest.NewRequest("GET", "/api/v1/concerts", nil)
	resp, err := ts.app.Test(req)
	if err != nil {
		t.Fatalf("failed to hit %s err %v", "/api/v1/concerts", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("response status code should be %d; got %d", fiber.StatusOK, resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	var apiResponse models.APIResponse[[]models.Concert]
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		t.Fatalf("could not unmarshal api response body err %v", err)
	}

	if apiResponse.Count != 1 {
		t.Fatalf("response count should be 1; got %d", apiResponse.Count)
	}
	if len(apiResponse.Data) != 1 {
		t.Fatalf("response data should have one concert; got %d", len(apiResponse.Data))
	}
	if apiResponse.Error != nil {
		t.Fatalf("response error should be nil; got %v", apiResponse.Error)
	}
	if apiResponse.Data[0].Id != concert2.Id {
		t.Fatalf("concert id does not match; want %d, got %d", concert2.Id, apiResponse.Data[0].Id)
	}
}
