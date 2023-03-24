package unit_test

import (
	"database/sql"
	"testing"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"
	"ticketbeastar/tests"
)

func TestTicketModel(t *testing.T) {
	db := database.OpenConnection(utils.GetTestDatabaseURL(), utils.InitLogger())
	defer database.CloseConnection(db)

	service := models.NewServices(db)

	testCases := map[string]func(t *testing.T){
		"can order concert tickets": func(t *testing.T) {
			concert := tests.CreateConcert(t, db, nil, true)
			email := "jane@example.com"
			var ticketQuanity uint64 = 3
			_, err := service.Ticket.Add(concert, ticketQuanity)
			if err != nil {
				t.Fatalf("could not create tickets; got %v", err)
			}

			order, err := service.Ticket.OrderTickets(email, concert.Id, ticketQuanity)
			if err != nil {
				t.Fatalf("could not order tickets; got %v", err)
			}

			count, err := service.Ticket.Remaining(concert)
			if err != nil {
				t.Fatalf("could not get remaining tickets; got %v", err)
			}
			if count != 0 {
				t.Fatalf("want %d tickets; got %d", ticketQuanity, count)
			}

			savedOrder, err := service.Order.FindByEmail(email)
			if err != nil {
				t.Fatalf("could not get order; %v", err)
			}
			if savedOrder.ConcertId != concert.Id {
				t.Fatalf("want order.conertId %d; got %d", savedOrder.ConcertId, concert.Id)
			}
			if savedOrder.Id != order.Id {
				t.Fatalf("want order %v; got order %v", savedOrder, order)
			}
		},
		"can add tickets to a concert": func(t *testing.T) {
			concert := tests.CreateConcert(t, db, nil, true)
			var ticketQuanity uint64 = 50

			tickets, err := service.Ticket.Add(concert, ticketQuanity)
			if err != nil {
				t.Fatalf("could not add tickets; got %v", err)
			}

			if len(*tickets) != int(ticketQuanity) {
				t.Fatalf("should have created %d tickets; got %d", ticketQuanity, len(*tickets))
			}

			count, err := service.Ticket.Remaining(concert)
			if err != nil {
				t.Fatalf("could not get remaining tickets; got %v", err)
			}
			if count != ticketQuanity {
				t.Fatalf("want %d tickets; got %d", ticketQuanity, count)
			}
			for i, ticket := range *tickets {
				if ticket.ConcertId != concert.Id {
					t.Fatalf("ticket(%d) should have concert id %d; got %d", i, concert.Id, ticket.ConcertId)
				}
			}
		},
		"tickets remaining does not include tickets associated with an order": func(t *testing.T) {
			concert := tests.CreateConcert(t, db, nil, true)
			var ticketQuanity uint64 = 50
			var orderedTicket uint64 = 30
			_, err := service.Ticket.Add(concert, ticketQuanity)
			if err != nil {
				t.Fatalf("could not create tickets; got %v", err)
			}

			_, err = service.Ticket.OrderTickets("jane@example.com", concert.Id, orderedTicket)
			if err != nil {
				t.Fatalf("could not order tickets; got %v", err)
			}

			count, err := service.Ticket.Remaining(concert)
			if err != nil {
				t.Fatalf("could not get remaining tickets; got %v", err)
			}
			if count != ticketQuanity-orderedTicket {
				t.Fatalf("want %d tickets remaining; got %d", ticketQuanity-orderedTicket, count)
			}
		},
		"trying to purchase more tickets than remaining returns an error": func(t *testing.T) {
			concert := tests.CreateConcert(t, db, nil, true)
			var ticketQuanity uint64 = 2
			email := "jane@example.com"
			_, err := service.Ticket.Add(concert, ticketQuanity)
			if err != nil {
				t.Fatalf("could not create tickets; got %v", err)
			}

			_, err = service.Ticket.OrderTickets(email, concert.Id, ticketQuanity+1)
			if err != models.ErrNotEnoughTickets {
				t.Fatalf("want ErrNotEnoughTickets; got %v", err)
			}

			if order, err := service.Order.FindByEmail(email); err != sql.ErrNoRows {
				t.Fatalf("order should be nil; got order(%v) err %v", order, err)
			}
			count, err := service.Ticket.Remaining(concert)
			if err != nil {
				t.Fatalf("could not get remaining tickets; got %v", err)
			}
			if count != ticketQuanity {
				t.Fatalf("want %d tickets remaining; got %d", ticketQuanity, count)
			}
		},
		"cannot order tickets that have already been purchased": func(t *testing.T) {
			concert := tests.CreateConcert(t, db, nil, true)
			var ticketQuanity uint64 = 10
			_, err := service.Ticket.Add(concert, ticketQuanity)
			if err != nil {
				t.Fatalf("could not create tickets; got %v", err)
			}

			_, err = service.Ticket.OrderTickets("jane@example.com", concert.Id, 8)
			if err != nil {
				t.Fatalf("could not order tickets; got %v", err)
			}
			_, err = service.Ticket.OrderTickets("micky@example.com", concert.Id, 3)
			if err != models.ErrNotEnoughTickets {
				t.Fatalf("want ErrNotEnoughTickets; got %v", err)
			}

			if mickeysOrder, err := service.Order.FindByEmail("micky@example.com"); err != sql.ErrNoRows {
				t.Fatalf("order should be nil; got order(%v) err %v", mickeysOrder, err)
			}
			count, err := service.Ticket.Remaining(concert)
			if err != nil {
				t.Fatalf("could not get remaining tickets; got %v", err)
			}
			if count != 2 {
				t.Fatalf("want %d tickets remaining; got %d", 2, count)
			}
		},
		"a ticket can be released": func(t *testing.T) {
			concert := tests.CreateConcert(t, db, nil, true)
			if _, err := service.Ticket.Add(concert, 1); err != nil {
				t.Fatalf("could not create tickets; got %v", err)
			}
			order, err := service.Ticket.OrderTickets("jane@example.com", concert.Id, 1)
			if err != nil {
				t.Fatalf("could not order tickets; got %v", err)
			}
			tickets, err := service.Ticket.FindByOrder(order.Id, 1)
			if err != nil {
				t.Fatalf("failed to fetch order tickets; %v", err)
			}
			if len(*tickets) != 1 {
				t.Fatalf("want %d ticket; got %d", 1, len(*tickets))
			}
			ticket := (*tickets)[0]
			if ticket.OrderId != order.Id {
				t.Fatalf("want ticket.orderId %d; %d", order.Id, ticket.OrderId)
			}

			service.Ticket.Release(order.Id)

			releasedTicket, err := service.Ticket.FindById(ticket.Id)
			if err != nil {
				t.Fatalf("failed to fetch ticket; %v", err)
			}
			if releasedTicket.OrderId != 0 {
				t.Fatalf("ticket should not have order id; got %d", releasedTicket.OrderId)
			}
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				tests.TeardownTicketTable(t, db)
				tests.TeardownOrderTable(t, db)
				tests.TeardownConcertTable(t, db)
			}()
			tests.SetupConcertTable(t, db)
			tests.SetupOrderTable(t, db)
			tests.SetupTicketable(t, db)

			tc(t)
		})
	}
}
