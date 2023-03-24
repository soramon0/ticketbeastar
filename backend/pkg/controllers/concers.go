package controllers

import (
	"database/sql"
	"errors"
	"log"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Concerts struct {
	concert models.ConcertService
	order   models.OrderService
	ticket  models.TicketService
	vt      *utils.ValidatorTransaltor
	log     *log.Logger
}

// New Users is used to create a new Users controller.
func NewConcerts(cs models.ConcertService, os models.OrderService, ts models.TicketService, vt *utils.ValidatorTransaltor, l *log.Logger) *Concerts {
	return &Concerts{
		concert: cs,
		order:   os,
		ticket:  ts,
		vt:      vt,
		log:     l,
	}
}

func (c *Concerts) GetConcerts(ctx *fiber.Ctx) error {
	concerts, err := c.concert.FindPublished()
	if err != nil {
		c.log.Println("GetConcerts", err)
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "internal server error"}
	}

	return ctx.JSON(models.NewAPIResponse(concerts, len(*concerts)))
}

func (c *Concerts) GetConcertById(ctx *fiber.Ctx) error {
	concert, err := c.concert.FindPublishedById(ctx.Params("id"))
	if err != nil {
		if err == models.ErrInvalidId {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: err.Error()}
		}
		if err == sql.ErrNoRows {
			return &fiber.Error{Code: fiber.StatusNotFound, Message: "Concert not found"}
		}

		c.log.Println("GetConcertById", err)
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "internal server error"}
	}

	return ctx.JSON(models.NewAPIResponse(concert, 0))
}

type CreateConcertOrderPayload struct {
	Email          string `json:"email" validate:"required,email,omitempty"`
	TicketQuantity int    `json:"ticket_quantity" validate:"required,number,gte=1,omitempty"`
	PaymentToken   string `json:"payment_token" validate:"required,omitempty"`
}

func (c *Concerts) CreateConcertOrder(ctx *fiber.Ctx) error {
	concert, err := c.concert.FindPublishedById(ctx.Params("id"))
	if err != nil {
		if err == models.ErrInvalidId {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: err.Error()}
		}
		if err == sql.ErrNoRows {
			return &fiber.Error{Code: fiber.StatusNotFound, Message: "Concert not found"}
		}

		c.log.Println("FindPublishedById", err)
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "internal server error"}
	}

	payload := new(CreateConcertOrderPayload)
	if err := ctx.BodyParser(payload); err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}
	if err := c.vt.Validator.Struct(payload); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return ctx.Status(fiber.StatusBadRequest).JSON(c.vt.ValidationErrors(ve))
		}
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: err.Error()}
	}

	order, err := c.ticket.OrderTickets(payload.Email, concert.Id, uint64(payload.TicketQuantity))
	if err != nil {
		if err == models.ErrNotEnoughTickets {
			return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
		}

		c.log.Println("Failed to create order", err)
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "internal server error"}
	}

	amount := uint64(payload.TicketQuantity) * concert.TicketPrice
	if err := chargePayment(amount, payload.PaymentToken); err != nil {
		if err := c.order.Cancel(order.Id); err != nil {
			c.log.Printf("failed to cancel order %q; got %v\n", order.Id, err)
		}

		if err == models.ErrInvalidPaymentToken {
			return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
		}

		c.log.Println("Failed to process payment", err)
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: err.Error()}
	}

	return ctx.Status(fiber.StatusCreated).JSON(models.NewAPIResponse(order, 0))
}

func chargePayment(amount uint64, token string) error {
	validPaymentToken := "valid payment token"
	if token != validPaymentToken {
		return models.ErrInvalidPaymentToken
	}

	return nil
}
