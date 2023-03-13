package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"ticketbeastar/pkg/models"

	"github.com/gofiber/fiber/v2"
)

type Concerts struct {
	service models.ConcertService
	log     *log.Logger
}

// New Users is used to create a new Users controller.
func NewConcerts(cs models.ConcertService, l *log.Logger) *Concerts {
	return &Concerts{
		service: cs,
		log:     l,
	}
}

func (c *Concerts) GetConcerts(ctx *fiber.Ctx) error {
	concerts, err := c.service.Find()
	if err != nil {
		c.log.Println("GetConcerts", err)
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "internal server error"}
	}

	return ctx.JSON(models.APIResponse{Data: concerts, Count: len(*concerts)})
}

func (c *Concerts) GetConcertById(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: fmt.Sprintf(`id %q is invalid`, ctx.Params("id"))}
	}

	concert, err := c.service.FindById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &fiber.Error{Code: fiber.StatusNotFound, Message: "Concert not found"}
		}

		c.log.Println("GetConcertById", err)
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "internal server error"}
	}

	return ctx.JSON(models.APIResponse{Data: concert})
}
