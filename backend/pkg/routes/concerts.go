package routes

import (
	"log"

	"ticketbeastar/pkg/controllers"
	"ticketbeastar/pkg/models"

	"github.com/gofiber/fiber/v2"
)

func registerConcertRoutes(a *fiber.App, s *models.Services, l *log.Logger) *fiber.Router {
	router := a.Group("/api/v1/concerts")
	concertsC := controllers.NewConcerts(s.Concert, l)

	router.Get("/", concertsC.GetConcerts)

	return &router
}
