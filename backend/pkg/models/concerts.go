package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Concert struct {
	bun.BaseModel `bun:"table:concerts,alias:c"`

	Id   int64  `bun:"id,pk,autoincrement" json:"id"`
	Name string `bun:"name,notnull" json:"name"`
}

type ConcertService interface {
	// Methods for querying users
	Find() (*[]Concert, error)

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

func (cs *concertService) Create(concert *Concert) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := cs.db.NewInsert().Model(concert).Exec(ctx)
	return err
}
