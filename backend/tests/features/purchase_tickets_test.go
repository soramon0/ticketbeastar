package features_test

import (
	"fmt"
	"net/http"
	"testing"
	"ticketbeastar/pkg/controllers"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/tests"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun/schema"
)

func TestPurchaseTickets(t *testing.T) {
	ts := tests.NewTestServer(t)
	defer database.CloseConnection(ts.Db)

	testsCases := map[string]func(t *testing.T){
		"customer can purchase concert ticket": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, &models.Concert{TicketPrice: 3250, PublishedAt: schema.NullTime{Time: time.Now()}}, true)

			endpoint := fmt.Sprintf("/api/v1/concerts/%d/orders", concert.Id)
			payload := controllers.CreateConcertOrderPayload{Email: "john@example.com", TicketPrice: 3, PaymentToken: "valid test token"}
			resp := ts.Json(t, http.MethodPost, endpoint, &payload)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusCreated)

			// assert order created for email
			// assert order for 3 tickets created
			// assert total price is ticket * quantity
			// ts.Service.Order.FindByEmail()
		},
	}

	for name, tc := range testsCases {
		t.Run(name, func(t *testing.T) {
			tests.SetupConcertTable(t, ts.Db)
			tests.SetupOrderTable(t, ts.Db)
			defer func() {
				tests.TeardownConcertTable(t, ts.Db)
				tests.TeardownOrderTable(t, ts.Db)
			}()
			tc(t)
		})
	}
}
