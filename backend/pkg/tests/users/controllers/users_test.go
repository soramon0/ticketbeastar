package controllers_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"
	"ticketbeastar/pkg/controllers"
	"ticketbeastar/pkg/models"

	"github.com/gofiber/fiber/v2"
)

func TestUsersController(t *testing.T) {
	ts := newTestServer()

	t.Run("List", func(t *testing.T) {
		ts.setup(t)
		defer ts.teardown(t)
		testUsersListing(t, ts)
	})

}

func testUsersListing(t *testing.T, ts *testServer) {
	users, err := createUsers(ts.db, 5)
	if err != nil {
		t.Fatalf("createUsers() err %v; want nil", err)
	}

	userC := controllers.NewUsers(ts.us, ts.log)
	ts.app.Get("/api/v1/users", userC.GetUsers)
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	resp, err := ts.app.Test(req)
	if err != nil {
		t.Fatal("GET /api/v1/users", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("statusCode = %d; want %d", resp.StatusCode, fiber.StatusOK)
	}
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var apiResponse models.APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		t.Fatalf("apiResponse unmarshal() err %v; want nil", err)
	}
	if apiResponse.Count != len(*users) {
		log.Fatalf("apiResponse.Count incorrect. want %d; got %d", apiResponse.Count, len(*users))
	}
	if apiResponse.Error != nil {
		log.Fatalf("apiResponse.Error got %v; want nil", apiResponse.Error)
	}
}
