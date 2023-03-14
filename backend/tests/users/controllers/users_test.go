package controllers_test

import (
	"encoding/json"
	"io"
	"testing"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/tests"
	"ticketbeastar/tests/users"

	"github.com/gofiber/fiber/v2"
)

func TestUsersController(t *testing.T) {
	ts := tests.NewTestServer()
	defer database.CloseConnection(ts.Db)

	tests := map[string]func(t *testing.T){
		"can list users": func(t *testing.T) {
			totalUsers := 5
			users := users.CreateUsers(t, ts.Db, uint(totalUsers))
			resp := ts.Visit(t, "/api/v1/users")
			api := unmarshalUsers(t, resp.Body)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusOK)
			ts.AssertResponseCount(t, api.Count, totalUsers)
			ts.AssertResponseError(t, api.Error, nil)

			data := *api.Data
			if len(data) != len(*users) {
				t.Fatalf("response data should have %d users; got %d", totalUsers, len(data))
			}
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			users.SetupTable(t, ts.Db)
			defer users.TeardownTable(t, ts.Db)
			test(t)
		})
	}
}

// func unmarshalUser(t *testing.T, body io.ReadCloser) models.APIResponse[*models.User] {
// 	content, err := io.ReadAll(body)
// 	if err != nil {
// 		t.Fatalf("could not read response body; err %v", err)
// 	}
// 	defer body.Close()

// 	var resp models.APIResponse[*models.User]
// 	if err := json.Unmarshal(content, &resp); err != nil {
// 		t.Fatalf("could not unmarshal concert response body; err %v", err)
// 	}
// 	return resp
// }

func unmarshalUsers(t *testing.T, body io.ReadCloser) models.APIResponse[*[]models.User] {
	content, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("could not read response body; err %v", err)
	}
	defer body.Close()

	var resp models.APIResponse[*[]models.User]
	if err := json.Unmarshal(content, &resp); err != nil {
		t.Fatalf("could not unmarshal concerts response body; err %v", err)
	}
	return resp
}
