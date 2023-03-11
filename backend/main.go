package main

import (
	"ticketbeastar/pkg/configs"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/routes"
	"ticketbeastar/pkg/utils"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger := utils.InitLogger()
	db := database.OpenConnection(utils.GetDatabaseURL(), logger)
	defer database.CloseConnection(db)

	app := fiber.New(configs.FiberConfig())
	services := models.NewServices(db)

	routes.Register(app, services, logger)
	utils.StartServer(app, logger)
}
