package routes

import (
	"log"

	"ticketbeastar/pkg/middleware"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func Register(a *fiber.App, s *models.Services, vt *utils.ValidatorTransaltor, l *log.Logger) {
	middleware.FiberMiddleware(a)

	registerUserRoutes(a, s, l)
	registerConcertRoutes(a, s, vt, l)
	registerNotFoundRoute(a)
}
