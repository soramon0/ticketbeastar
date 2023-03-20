package features_test

import (
	"fmt"
	"net/http"
	"testing"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/tests"

	"github.com/gofiber/fiber/v2"
)

type orderPayload struct {
	Email           string
	Ticket_quantity int32
	Payment_token   string
}

func TestPurchaseTickets(t *testing.T) {
	ts := tests.NewTestServer(t)
	defer database.CloseConnection(ts.Db)

	testsCases := map[string]func(t *testing.T){
		"customer can purchase concert ticket": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, &models.Concert{TicketPrice: 3250}, true)

			endpoint := fmt.Sprintf("/api/v1/concerts/%d/orders", concert.Id)
			op := orderPayload{Email: "john@example.com", Ticket_quantity: 3, Payment_token: ""}
			resp := ts.Json(t, http.MethodPost, endpoint, op)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusCreated)
		},
	}

	for name, tc := range testsCases {
		t.Run(name, func(t *testing.T) {
			tests.SetupConcertTable(t, ts.Db)
			defer tests.TeardownConcertTable(t, ts.Db)
			tc(t)
		})
	}
}
