package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders,alias:o"`

	Id        uint64    `bun:"id,pk,autoincrement" json:"id"`
	Email     string    `bun:"email,notnull" json:"email"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`

	ConcertId uint64    `bun:",notnull" json:"concert_id"`
	Tickets   []*Ticket `bun:"rel:has-many,join:id=order_id" json:"tickets"`
}

type OrderService interface {
	// Methods for querying orders
	Find() (*[]Order, error)
	FindById(id uint64) (*Order, error)
	FindByEmail(email string) (*Order, error)

	// Creates a concert order for the given user email
	Create(email string, concertId uint64) (*Order, error)
	// Will release all tickets for the order and then delete the order
	Cancel(orderId uint64) error
	// Deletes the order with the given id
	Delete(orderId uint64) error
}

type orderService struct {
	db *bun.DB
}

func NewOrderService(db *bun.DB) OrderService {
	return &orderService{
		db: db,
	}
}

func (os *orderService) Find() (*[]Order, error) {
	orders := []Order{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := os.db.NewSelect().Model(&orders).Scan(ctx)
	return &orders, err
}

func (os *orderService) FindById(id uint64) (*Order, error) {
	var order Order
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := os.db.NewSelect().Model(&order).Where("id = ?", id).Scan(ctx)
	return &order, err
}

func (os *orderService) FindByEmail(email string) (*Order, error) {
	var order Order
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := os.db.NewSelect().Model(&order).Where("email = ?", email).Scan(ctx)
	return &order, err
}

func (os *orderService) Create(email string, concertId uint64) (*Order, error) {
	return createOrder(os.db, email, concertId)
}

func createOrder(db *bun.DB, email string, concertId uint64) (*Order, error) {
	order := &Order{Email: email, ConcertId: concertId}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := db.NewInsert().Model(order).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (os *orderService) Cancel(orderId uint64) error {
	if err := releaseTickets(os.db, orderId); err != nil {
		return err
	}
	return os.Delete(orderId)
}

func (os *orderService) Delete(orderId uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := os.db.NewDelete().Model((*Order)(nil)).Where("id = ?", orderId).Exec(ctx)
	return err
}
