package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Concert struct {
	bun.BaseModel `bun:"table:concerts,alias:c"`

	Id                    uint64    `bun:"id,pk,autoincrement" json:"id"`
	Title                 string    `bun:"title,notnull" json:"title"`
	Subtitle              string    `bun:"subtitle,notnull" json:"subtitle"`
	Date                  time.Time `bun:"date,notnull" json:"date"`
	TicketPrice           uint64    `bun:"ticket_price,notnull" json:"ticket_price"`
	Venue                 string    `bun:"venue,notnull" json:"venue"`
	VenueAddress          string    `bun:"venue_address,notnull" json:"venue_address"`
	City                  string    `bun:"city,notnull" json:"city"`
	State                 string    `bun:"state,notnull" json:"state"`
	Zip                   string    `bun:"zip,notnull" json:"zip"`
	AdditionalInformation string    `bun:"additional_information,type:text,notnull" json:"additional_information"`
	CreatedAt             time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt             time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

type ConcertService interface {
	// Methods for querying users
	Find() (*[]Concert, error)
	FindById(id uint64) (*Concert, error)

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	concerts := []Concert{}
	err := cs.db.NewSelect().Model(&concerts).Scan(ctx)
	return &concerts, err
}

func (cs *concertService) FindById(id uint64) (*Concert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var concert Concert
	err := cs.db.NewSelect().Model(&concert).Where("id = ?", id).Scan(ctx)
	return &concert, err
}

func (cs *concertService) Create(concert *Concert) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := cs.db.NewInsert().Model(concert).Exec(ctx)
	return err
}
