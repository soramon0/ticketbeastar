package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Ticket struct {
	bun.BaseModel `bun:"table:tickets,alias:t"`

	Id        uint64    `bun:"id,pk,autoincrement" json:"id"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`

	OrderId uint64 `json:"order_id,omitempty"`
}

type TicketService interface {
	// Methods for querying tickets
	Find() (*[]Ticket, error)
	// FindById(id uint64) (*Ticket, error)
	// FindByEmail(email string) (*Ticket, error)

	// Methods for altering tickets
	Create(ticket *Ticket) error
	BulkCreate(tickets *[]Ticket) error
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
