package features_test

import (
	"encoding/json"
	"fmt"
	"io"
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
			email := "john@example.com"
			var ticketQuantity uint64 = 3
			payload := controllers.CreateConcertOrderPayload{Email: email, TicketQuantity: ticketQuantity, PaymentToken: "valid test token"}
			resp := ts.Json(t, http.MethodPost, endpoint, &payload)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusCreated)

			order, err := ts.Service.Order.FindByEmail(email)
			if err != nil {
				t.Fatalf("order should not be %v; %v", order, err)
			}
			if order.ConcertId != concert.Id {
				t.Fatalf("order concert id should be %d; got %d", order.ConcertId, concert.Id)
			}

			tickets, err := ts.Service.Ticket.Find()
			if err != nil {
				t.Fatalf("tickets should not be %v; %v", order, err)
			}
			if len(*tickets) != int(ticketQuantity) {
				t.Fatalf("should have created %d tickets; got %d", ticketQuantity, len(*tickets))
			}
			for i, ticket := range *tickets {
				if ticket.OrderId != order.Id {
					t.Fatalf("ticket(%d) should have order id %d; got %d", i, order.Id, ticket.OrderId)
				}
			}
			// assert total price is ticket * quantity
		},
		"email is required to purchase tickets": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, nil, true)
			endpoint := fmt.Sprintf("/api/v1/concerts/%d/orders", concert.Id)
			payload := controllers.CreateConcertOrderPayload{Email: "", TicketQuantity: 3, PaymentToken: "valid test token"}
			resp := ts.Json(t, http.MethodPost, endpoint, &payload)
			gotErr := unmarshalValidationErrors(t, resp.Body)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusBadRequest)

			if len(gotErr.Errors) != 1 {
				t.Fatalf("want %d validation error(s); got %d", 1, len(gotErr.Errors))
			}

			wantErr := &models.APIValidaitonErrors{Errors: []models.APIFieldError{
				{Field: "email", Message: "email is required"},
			}}
			ts.AssertResponseValidationError(t, gotErr, wantErr)
		},
	}

	for name, tc := range testsCases {
		t.Run(name, func(t *testing.T) {
			tests.SetupConcertTable(t, ts.Db)
			tests.SetupOrderTable(t, ts.Db)
			tests.SetupTicketable(t, ts.Db)
			defer func() {
				tests.TeardownTicketTable(t, ts.Db)
				tests.TeardownOrderTable(t, ts.Db)
				tests.TeardownConcertTable(t, ts.Db)
			}()
			tc(t)
		})
	}
}

func unmarshalValidationErrors(t *testing.T, body io.ReadCloser) *models.APIValidaitonErrors {
	content, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("could not read response body; err %v", err)
	}
	defer body.Close()

	var resp models.APIValidaitonErrors
	if err := json.Unmarshal(content, &resp); err != nil {
		t.Fatalf("could not unmarshal validation errors; err %v", err)
	}
	return &resp
}
