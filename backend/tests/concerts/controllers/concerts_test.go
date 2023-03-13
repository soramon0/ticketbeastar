package controllers_test

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestConcertsController(t *testing.T) {
	ts := newTestServer()

	t.Run("List single concert", func(t *testing.T) {
		ts.setup(t)
		defer ts.teardown(t)
		testConcertListing(t, ts)
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
