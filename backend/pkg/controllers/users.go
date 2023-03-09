package controllers

import (
	"log"

	"ticketbeastar/pkg/models"

	"github.com/gofiber/fiber/v2"
)

type Users struct {
	service models.UserService
	log     *log.Logger
}

// New Users is used to create a new Users controller.
func NewUsers(us models.UserService, l *log.Logger) *Users {
	return &Users{
		service: us,
		log:     l,
	}
}

func (u *Users) GetUsers(c *fiber.Ctx) error {
	users, err := u.service.Find()
	if err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: err.Error()}
	}

	return c.JSON(models.APIResponse{Data: users, Count: len(*users)})
}
