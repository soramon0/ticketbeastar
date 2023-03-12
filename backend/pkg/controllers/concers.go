package controllers

import (
	"log"
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
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	return ctx.JSON(models.APIResponse{Data: concerts, Count: len(*concerts)})
}
