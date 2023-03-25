package features_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"ticketbeastar/pkg/controllers"
	"ticketbeastar/pkg/models"
	"ticketbeastar/tests"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

func TestPurchaseTickets(t *testing.T) {
	ts := tests.NewTestServer(t)
	defer ts.CloseConnection(t)

	validPaymentToken := "valid payment token"

	testsCases := map[string]func(t *testing.T){
		"customer can purchase tickets to a published concert": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, &models.Concert{TicketPrice: 3250, PublishedAt: schema.NullTime{Time: time.Now()}}, true)
			email := "john@example.com"
			var ticketQuantity uint64 = 3
			ts.Service.Ticket.Add(concert, ticketQuantity)
			payload := controllers.CreateConcertOrderPayload{Email: email, TicketQuantity: int(ticketQuantity), PaymentToken: validPaymentToken}
			resp := orderTickets(t, ts, concert.Id, payload)
			api := unmarshalOrder(t, resp.Body)

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

			ts.AssertResponseError(t, api.Error, nil)
			assertOrder(t, api.Data, email, ticketQuantity, 3250*ticketQuantity)

			// assert no charge was made
		},
		"customer cannot purchase tickets to an unpublished concert": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, &models.Concert{PublishedAt: bun.NullTime{}}, true)
			ts.Service.Ticket.Add(concert, 3)
			email := "john@example.com"
			payload := controllers.CreateConcertOrderPayload{Email: email, TicketQuantity: 1, PaymentToken: validPaymentToken}
			resp := orderTickets(t, ts, concert.Id, payload)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusNotFound)

			api := tests.UnmarshalConcert(t, resp.Body)
			ts.AssertResponseError(t, api.Error, &models.APIError{Message: "Concert not found"})

			if _, err := ts.Service.Order.FindByEmail(email); err != sql.ErrNoRows {
				t.Fatalf("no order should be created; got %v", err)
			}
			// assert no charge was made
		},
		"cannot purchase more tickets than available": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, nil, true)
			ts.Service.Ticket.Add(concert, 50)

			email := "john@example.com"
			payload := controllers.CreateConcertOrderPayload{Email: email, TicketQuantity: 51, PaymentToken: validPaymentToken}
			resp := orderTickets(t, ts, concert.Id, payload)
			api := unmarshalOrder(t, resp.Body)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusUnprocessableEntity)
			ts.AssertResponseError(t, api.Error, &models.APIError{Message: models.ErrNotEnoughTickets.Error()})

			if _, err := ts.Service.Order.FindByEmail(email); err != sql.ErrNoRows {
				t.Fatalf("no order should be created; got %v", err)
			}

			count, err := ts.Service.Ticket.Remaining(concert)
			if err != nil {
				t.Fatalf("could not get remaining tickets; got %v", err)
			}
			if count != 50 {
				t.Fatalf("want %d tickets remaining; got %d", 50, count)
			}
			// assert no charge was made
		},
		"email is required to purchase tickets": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, nil, true)
			payload := controllers.CreateConcertOrderPayload{Email: "", TicketQuantity: 3, PaymentToken: validPaymentToken}
			resp := orderTickets(t, ts, concert.Id, payload)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusBadRequest)

			gotErr := unmarshalValidationErrors(t, resp.Body)
			if len(gotErr.Errors) != 1 {
				t.Fatalf("want %d validation error(s); got %d", 1, len(gotErr.Errors))
			}
			wantErr := &models.APIValidaitonErrors{Errors: []models.APIFieldError{
				{Field: "email", Message: "email is required"},
			}}
			ts.AssertResponseValidationError(t, gotErr, wantErr)
		},
		"email must be valid to purchase tickets": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, nil, true)
			payload := controllers.CreateConcertOrderPayload{Email: "not-a-email-address", TicketQuantity: 3, PaymentToken: validPaymentToken}
			resp := orderTickets(t, ts, concert.Id, payload)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusBadRequest)

			gotErr := unmarshalValidationErrors(t, resp.Body)
			if len(gotErr.Errors) != 1 {
				t.Fatalf("want %d validation error(s); got %d", 1, len(gotErr.Errors))
			}
			wantErr := &models.APIValidaitonErrors{Errors: []models.APIFieldError{
				{Field: "email", Message: "email must be a valid email address"},
			}}
			ts.AssertResponseValidationError(t, gotErr, wantErr)
		},
		"ticket_quantity is required to purchase tickets": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, nil, true)
			payload := controllers.CreateConcertOrderPayload{Email: "jon@example.com", TicketQuantity: 0, PaymentToken: validPaymentToken}
			resp := orderTickets(t, ts, concert.Id, payload)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusBadRequest)

			gotErr := unmarshalValidationErrors(t, resp.Body)
			if len(gotErr.Errors) != 1 {
				t.Fatalf("want %d validation error(s); got %d", 1, len(gotErr.Errors))
			}
			wantErr := &models.APIValidaitonErrors{Errors: []models.APIFieldError{
				{Field: "ticket_quantity", Message: "ticket_quantity is required"},
			}}
			ts.AssertResponseValidationError(t, gotErr, wantErr)
		},
		"ticket_quantity must at least be 1 to purchase tickets": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, nil, true)
			payload := controllers.CreateConcertOrderPayload{Email: "jon@example.com", TicketQuantity: -1, PaymentToken: validPaymentToken}
			resp := orderTickets(t, ts, concert.Id, payload)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusBadRequest)

			gotErr := unmarshalValidationErrors(t, resp.Body)
			if len(gotErr.Errors) != 1 {
				t.Fatalf("want %d validation error(s); got %d", 1, len(gotErr.Errors))
			}
			wantErr := &models.APIValidaitonErrors{Errors: []models.APIFieldError{
				{Field: "ticket_quantity", Message: "ticket_quantity must be 1 or greater"},
			}}
			ts.AssertResponseValidationError(t, gotErr, wantErr)
		},
		"payment_token is required to purchase tickets": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, nil, true)
			payload := controllers.CreateConcertOrderPayload{Email: "jon@example.com", TicketQuantity: 1, PaymentToken: ""}
			resp := orderTickets(t, ts, concert.Id, payload)

			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusBadRequest)

			gotErr := unmarshalValidationErrors(t, resp.Body)
			if len(gotErr.Errors) != 1 {
				t.Fatalf("want %d validation error(s); got %d", 1, len(gotErr.Errors))
			}
			wantErr := &models.APIValidaitonErrors{Errors: []models.APIFieldError{
				{Field: "payment_token", Message: "payment_token is required"},
			}}
			ts.AssertResponseValidationError(t, gotErr, wantErr)
		},
		"an order is not created if payment fails": func(t *testing.T) {
			concert := tests.CreateConcert(t, ts.Db, nil, true)
			ts.Service.Ticket.Add(concert, 1)
			email := "jon@example.com"
			payload := controllers.CreateConcertOrderPayload{Email: email, TicketQuantity: 1, PaymentToken: "invalid payment token"}
			resp := orderTickets(t, ts, concert.Id, payload)
			ts.AssertResponseStatus(t, resp.StatusCode, fiber.StatusUnprocessableEntity)

			api := unmarshalOrder(t, resp.Body)
			ts.AssertResponseError(t, api.Error, &models.APIError{Message: models.ErrInvalidPaymentToken.Error()})

			_, err := ts.Service.Order.FindByEmail(email)
			if err != sql.ErrNoRows {
				t.Fatalf("no order should be created; got %v", err)
			}
		},
	}

	for name, tc := range testsCases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				tests.TeardownTicketTable(t, ts.Db)
				tests.TeardownOrderTable(t, ts.Db)
				tests.TeardownConcertTable(t, ts.Db)
			}()
			tests.SetupConcertTable(t, ts.Db)
			tests.SetupOrderTable(t, ts.Db)
			tests.SetupTicketable(t, ts.Db)
			tc(t)
		})
	}
}

func orderTickets(t *testing.T, ts *tests.TestServer, concertId uint64, payload controllers.CreateConcertOrderPayload) *http.Response {
	endpoint := fmt.Sprintf("/api/v1/concerts/%d/orders", concertId)
	return ts.Json(t, http.MethodPost, endpoint, payload)
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

func unmarshalOrder(t *testing.T, body io.ReadCloser) models.APIResponse[*models.Order] {
	content, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("could not read response body; err %v", err)
	}
	defer body.Close()

	var resp models.APIResponse[*models.Order]
	if err := json.Unmarshal(content, &resp); err != nil {
		t.Fatalf("could not unmarshal concerts response body; err %v", err)
	}
	return resp
}

func assertOrder(t *testing.T, order *models.Order, email string, ticketQuantity uint64, amount uint64) {
	if order.Email != email {
		t.Fatalf("want email %s, got %s", email, order.Email)
	}
	if order.TicketQuantity != ticketQuantity {
		t.Fatalf("want ticketQuantity %d, got %d", ticketQuantity, order.TicketQuantity)
	}
	if order.Amount != amount {
		t.Fatalf("want amount %d, got %d", amount, order.Amount)
	}
}
