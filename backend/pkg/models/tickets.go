package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type Ticket struct {
	bun.BaseModel `bun:"table:tickets,alias:t"`

	Id        uint64    `bun:"id,pk,autoincrement" json:"id"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`

	ConcertId uint64 `bun:",notnull" json:"concert_id,omitempty"`
	OrderId   uint64 `bun:",nullzero" json:"order_id,omitempty"`
}

type TicketService interface {
	// Methods for querying tickets
	Find() (*[]Ticket, error)
	FindByConcert(concertId uint64, limit int) (*[]Ticket, error)
	// FindById(id uint64) (*Ticket, error)
	// FindByEmail(email string) (*Ticket, error)

	// Methods for altering tickets
	Create(ticket *Ticket) error
	BulkCreate(tickets *[]Ticket) error
	// Uses ticketQuantity to create tickets for an order
	CreateOrderTickets(order *Order, ticketQuantity uint64) (*[]Ticket, error)
	// Uses ticketQuantity to generate tickets for a concert
	Add(concert *Concert, ticketQuantity uint64) (*[]Ticket, error)
	// returns the number of remaining tickets for concert
	Remaining(concert *Concert) (uint64, error)
}

type ticketService struct {
	db *bun.DB
}

func NewTicketService(db *bun.DB) TicketService {
	return &ticketService{
		db: db,
	}
}

func (os *ticketService) Find() (*[]Ticket, error) {
	tickets := []Ticket{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := os.db.NewSelect().Model(&tickets).Scan(ctx)
	return &tickets, err
}

func (os *ticketService) FindByConcert(concertId uint64, limit int) (*[]Ticket, error) {
	tickets := []Ticket{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := os.db.NewSelect().Model(&tickets).Where("concert_id = ?", concertId).Limit(limit).Scan(ctx)
	if err != nil {
		return nil, err
	}
	if len(tickets) == 0 {
		return nil, sql.ErrNoRows
	}

	return &tickets, err
}

func (os *ticketService) Create(ticket *Ticket) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := os.db.NewInsert().Model(ticket).Exec(ctx)
	return err
}

func (os *ticketService) BulkCreate(tickets *[]Ticket) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := os.db.NewInsert().Model(tickets).Exec(ctx)
	return err
}

func (os *ticketService) CreateOrderTickets(order *Order, ticketQuantity uint64) (*[]Ticket, error) {
	tickets, err := os.FindByConcert(order.ConcertId, int(ticketQuantity))
	if err != nil {
		return nil, err
	}

	for i := range *tickets {
		(*tickets)[i].OrderId = order.Id
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = os.db.NewUpdate().Model(tickets).Column("order_id").Bulk().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (os *ticketService) Add(concert *Concert, ticketQuantity uint64) (*[]Ticket, error) {
	tickets := make([]Ticket, ticketQuantity)
	for i := range tickets {
		tickets[i].ConcertId = concert.Id
		concert.Tickets = append(concert.Tickets, &tickets[i])
	}

	if err := os.BulkCreate(&tickets); err != nil {
		return nil, err
	}
	return &tickets, nil
}

func (os *ticketService) Remaining(concert *Concert) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := os.db.NewSelect().Model((*Ticket)(nil)).Where("order_id IS NULL").Count(ctx)
	return uint64(count), err
}
