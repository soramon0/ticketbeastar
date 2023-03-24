package unit_test

import (
	"database/sql"
	"testing"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"
	"ticketbeastar/tests"
)

func TestOrderModel(t *testing.T) {
	db := database.OpenConnection(utils.GetTestDatabaseURL(), false, utils.InitLogger())
	defer database.CloseConnection(db)

	service := models.NewServices(db)

	testCases := map[string]func(t *testing.T, cs models.ConcertService){
		"tickets are released when an order is cancelled": func(t *testing.T, cs models.ConcertService) {
			concert := tests.CreateConcert(t, db, nil, true)
			if _, err := service.Ticket.Add(concert, 10); err != nil {
				t.Fatalf("could not create tickets; got %v", err)
			}
			order, err := service.Ticket.OrderTickets(concert, "jane@example.com", 5)
			if err != nil {
				t.Fatalf("could not order tickets; got %v", err)
			}
			count, err := service.Ticket.Remaining(concert)
			if err != nil {
				t.Fatalf("could not get remaining tickets; got %v", err)
			}
			if count != 5 {
				t.Fatalf("want %d tickets remaining; got %d", 5, count)
			}

			if err = service.Order.Cancel(order.Id); err != nil {
				t.Fatalf("could not cancel order; %v", err)
			}

			count, err = service.Ticket.Remaining(concert)
			if err != nil {
				t.Fatalf("could not get remaining tickets; got %v", err)
			}
			if count != 10 {
				t.Fatalf("want %d tickets remaining; got %d", 5, count)
			}
			if order, err := service.Order.FindById(order.Id); err != sql.ErrNoRows {
				t.Fatalf("order should be nil; got order(%v) err %v", order, err)
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

			tc(t, models.NewConcertService(db))
		})
	}
}
