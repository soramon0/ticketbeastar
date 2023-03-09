package routes

import "github.com/gofiber/fiber/v2"

func registerNotFoundRoute(a *fiber.App) {
	a.Use(
		func(c *fiber.Ctx) error {
			return &fiber.Error{Code: fiber.StatusNotFound, Message: "sorry, endpoint is not found"}
		},
	)
}
