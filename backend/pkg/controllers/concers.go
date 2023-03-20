package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Concerts struct {
	service models.ConcertService
	vt      *utils.ValidatorTransaltor
	log     *log.Logger
}

// New Users is used to create a new Users controller.
func NewConcerts(cs models.ConcertService, vt *utils.ValidatorTransaltor, l *log.Logger) *Concerts {
	return &Concerts{
		service: cs,
		vt:      vt,
		log:     l,
	}
}

func (c *Concerts) GetConcerts(ctx *fiber.Ctx) error {
	concerts, err := c.service.FindPublished()
	if err != nil {
		c.log.Println("GetConcerts", err)
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "internal server error"}
	}

	return ctx.JSON(models.NewAPIResponse(concerts, len(*concerts), nil))
}

func (c *Concerts) GetConcertById(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: fmt.Sprintf(`id %q is invalid`, ctx.Params("id"))}
	}

	concert, err := c.service.FindPublishedById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &fiber.Error{Code: fiber.StatusNotFound, Message: "Concert not found"}
		}

		c.log.Println("GetConcertById", err)
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "internal server error"}
	}

	return ctx.JSON(models.NewAPIResponse(concert, 0, nil))
}

type CreateConcertOrderPayload struct {
	Email        string `json:"email" validate:"required,email,omitempty"`
	TicketPrice  int64  `json:"ticket_price" validate:"required,number,gte=0,omitempty"`
	PaymentToken string `json:"payment_token" validate:"required,omitempty"`
}

func (c *Concerts) CreateConcertOrder(ctx *fiber.Ctx) error {
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

	return ctx.SendStatus(fiber.StatusCreated)
}
