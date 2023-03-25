package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Concert struct {
	bun.BaseModel `bun:"table:concerts,alias:c"`

	Id                    uint64       `bun:"id,pk,autoincrement" json:"id"`
	Title                 string       `bun:"title,notnull" json:"title"`
	Subtitle              string       `bun:"subtitle,notnull" json:"subtitle"`
	Date                  time.Time    `bun:"date,notnull" json:"date"`
	TicketPrice           uint64       `bun:"ticket_price,notnull" json:"ticket_price"`
	Venue                 string       `bun:"venue,notnull" json:"venue"`
	VenueAddress          string       `bun:"venue_address,notnull" json:"venue_address"`
	City                  string       `bun:"city,notnull" json:"city"`
	State                 string       `bun:"state,notnull" json:"state"`
	Zip                   string       `bun:"zip,notnull" json:"zip"`
	AdditionalInformation string       `bun:"additional_information,type:text,notnull" json:"additional_information"`
	PublishedAt           bun.NullTime `bun:"published_at,nullzero" json:"published_at"`
	CreatedAt             time.Time    `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt             time.Time    `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Orders  []*Order  `bun:"rel:has-many,join:id=concert_id" json:"orders,omitempty"`
	Tickets []*Ticket `bun:"rel:has-many,join:id=concert_id" json:"tickets,omitempty"`
}

type ConcertService interface {
	// Methods for querying orders
	Find() (*[]Concert, error)
	FindPublished() (*[]Concert, error)
	FindById(id string) (*Concert, error)
	FindPublishedById(id string) (*Concert, error)
	FindTickets(concertId uint64, ticketQuantity uint64) (*[]Ticket, error)
	CreateOrder(email string, concert *Concert, tickets *[]Ticket) (*Order, error)

	// Methods for altering concerts
	Create(concert *Concert) error
}

type concertService struct {
	db *bun.DB
}

func NewConcertService(db *bun.DB) ConcertService {
	return &concertService{
		db: db,
	}
}

func (cs *concertService) Find() (*[]Concert, error) {
	concerts := []Concert{}
	query := buildSelectQuery(cs.db, &concerts, false)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return &concerts, query.Scan(ctx)
}

func (cs *concertService) FindPublished() (*[]Concert, error) {
	concerts := []Concert{}
	query := buildSelectQuery(cs.db, &concerts, true)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return &concerts, query.Scan(ctx)
}

func (cs *concertService) FindById(id string) (*Concert, error) {
	parsedId, err := parseId(id)
	if err != nil {
		return nil, err
	}

	var concert Concert
	query := buildSelectQuery(cs.db, &concert, false).Where("id = ?", parsedId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return &concert, query.Scan(ctx)
}

func (cs *concertService) FindPublishedById(id string) (*Concert, error) {
	parsedId, err := parseId(id)
	if err != nil {
		return nil, err
	}

	var concert Concert
	query := buildSelectQuery(cs.db, &concert, true).Where("id = ?", parsedId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return &concert, query.Scan(ctx)
}

func (cs *concertService) FindTickets(concertId uint64, ticketQuantity uint64) (*[]Ticket, error) {
	tickets, err := findTicketsForConcert(cs.db, concertId, int(ticketQuantity))

	if len(*tickets) < int(ticketQuantity) {
		return nil, ErrNotEnoughTickets
	}

	return tickets, err
}

func (cs *concertService) CreateOrder(email string, concert *Concert, tickets *[]Ticket) (*Order, error) {
	order, err := createOrder(cs.db, concert, email, uint64(len(*tickets)))
	if err != nil {
		return nil, err
	}

	for i := range *tickets {
		(*tickets)[i].OrderId = order.Id
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = cs.db.NewUpdate().Model(tickets).Column("order_id").Bulk().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (cs *concertService) Create(concert *Concert) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := cs.db.NewInsert().Model(concert).Exec(ctx)
	return err
}

func buildSelectQuery(db *bun.DB, model any, published bool) *bun.SelectQuery {
	query := db.NewSelect().Model(model)
	if published {
		query.Where("published_at IS NOT NULL")
	}
	return query
}
