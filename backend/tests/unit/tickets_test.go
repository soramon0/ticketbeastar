package unit_test

import (
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

	testCases := map[string]func(t *testing.T, cs models.ConcertService){
		"can order concert tickets": func(t *testing.T, cs models.ConcertService) {
			concert := tests.CreateConcert(t, db, nil, true)
			email := "jane@example.com"
			var ticketQuanity uint64 = 3
			order, err := service.Order.Create(email, concert.Id)
			if err != nil {
				t.Fatalf("could not create order; got %v", err)
			}
			_, err = service.Ticket.Add(concert, ticketQuanity)
			if err != nil {
				t.Fatalf("could not create tickets; got %v", err)
			}
			tickets, err := service.Ticket.CreateOrderTickets(order, ticketQuanity)
			if err != nil {
				t.Fatalf("could not create order tickets; got %v", err)
			}

			if len(*tickets) != int(ticketQuanity) {
				t.Fatalf("should have created %d tickets; got %d", ticketQuanity, len(*tickets))
			}

			for i, ticket := range *tickets {
				if ticket.OrderId != order.Id {
					t.Fatalf("ticket(%d) should have order id %d; got %d", i, order.Id, ticket.OrderId)
				}
			}
		},
		"can add tickets to a concert": func(t *testing.T, cs models.ConcertService) {
			concert := tests.CreateConcert(t, db, nil, true)
			var ticketQuanity uint64 = 50

			tickets, err := service.Ticket.Add(concert, ticketQuanity)
			if err != nil {
				t.Fatalf("could not create tickets; got %v", err)
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
		"tickets remaining does not include tickets associated with an order": func(t *testing.T, cs models.ConcertService) {
			concert := tests.CreateConcert(t, db, nil, true)
			var ticketQuanity uint64 = 50
			var orderedTicket uint64 = 30
			_, err := service.Ticket.Add(concert, ticketQuanity)
			if err != nil {
				t.Fatalf("could not create tickets; got %v", err)
			}
			order, err := service.Order.Create("jane@example.com", concert.Id)
			if err != nil {
				t.Fatalf("could not create order; got %v", err)
			}
			_, err = service.Ticket.CreateOrderTickets(order, orderedTicket)
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
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tests.SetupConcertTable(t, db)
			tests.SetupOrderTable(t, db)
			tests.SetupTicketable(t, db)
			defer func() {
				tests.TeardownTicketTable(t, db)
				tests.TeardownOrderTable(t, db)
				tests.TeardownConcertTable(t, db)
			}()

			tc(t, models.NewConcertService(db))
		})
	}
}
