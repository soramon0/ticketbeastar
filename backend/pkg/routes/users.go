package routes

import (
	"log"

	"ticketbeastar/pkg/controllers"
	"ticketbeastar/pkg/models"

	"github.com/gofiber/fiber/v2"
)

func registerUserRoutes(a *fiber.App, s *models.Services, l *log.Logger) *fiber.Router {
	router := a.Group("/api/v1/users")
	usersC := controllers.NewUsers(s.User, l)

	router.Get("/", usersC.GetUsers)

	return &router
}
