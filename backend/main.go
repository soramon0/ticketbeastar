package main

import (
	"flag"
	"ticketbeastar/pkg/configs"
	"ticketbeastar/pkg/database"
	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/routes"
	"ticketbeastar/pkg/utils"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var verboseDb bool
	flag.BoolVar(&verboseDb, "db-verbose", false, "Used to enable db verbose logs")
	flag.Parse()

	logger := utils.InitLogger()
	db := database.OpenConnection(utils.GetDatabaseURL(), verboseDb, logger)
	defer func() {
		utils.Must(database.CloseConnection(db))
	}()

	app := fiber.New(configs.FiberConfig())
	services := models.NewServices(db)

	vt, err := utils.NewValidator()
	if err != nil {
		logger.Fatalf("could not create validator %v", err)
	}

	routes.Register(app, services, vt, logger)
	utils.StartServer(app, logger)
}
