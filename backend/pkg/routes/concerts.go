package routes

import (
	"log"

	"ticketbeastar/pkg/controllers"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func registerConcertRoutes(a *fiber.App, s *models.Services, vt *utils.ValidatorTransaltor, l *log.Logger) *fiber.Router {
	router := a.Group("/api/v1/concerts")
	concertsC := controllers.NewConcerts(s.Concert, s.Order, vt, l)

	router.Get("/", concertsC.GetConcerts)
	router.Get("/:id", concertsC.GetConcertById)
	router.Post("/:id/orders", concertsC.CreateConcertOrder)

	return &router
}
