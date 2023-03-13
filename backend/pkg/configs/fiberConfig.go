package configs

import (
	"errors"
	"strconv"
	"time"

	"ticketbeastar/pkg/models"
	"ticketbeastar/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// FiberConfig func for configuration Fiber app.
// See: https://docs.gofiber.io/api/fiber#config
func FiberConfig() fiber.Config {
	// Define server settings.
	readTimeoutSecondsCount, _ := strconv.Atoi(utils.GetServerReadTimeout())

	// Return Fiber configuration.
	return fiber.Config{
		ReadTimeout:  time.Second * time.Duration(readTimeoutSecondsCount),
		ErrorHandler: errHandler,
	}
}

func errHandler(ctx *fiber.Ctx, err error) error {
	// Status code and message defaults to 500
	apiError := models.APIResponse[any]{Error: &models.APIError{
		Message:    fiber.ErrInternalServerError.Error(),
		StatusCode: fiber.StatusInternalServerError,
	}}

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		apiError.Error.StatusCode = e.Code
		apiError.Error.Message = e.Message
	}

	// Send custom error response
	if err := ctx.Status(apiError.Error.StatusCode).JSON(apiError); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[any]{
			Error: &models.APIError{Message: "Internal Server Error"},
		})
	}

	// Return from handler
	return nil
}
