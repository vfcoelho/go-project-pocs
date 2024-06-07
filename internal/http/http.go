package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	internal "github.com/vfcoelho/go-project-pocs/internal"
	errs "github.com/vfcoelho/go-project-pocs/internal/errors"
	"github.com/vfcoelho/go-project-pocs/src/dtos"
	"github.com/vfcoelho/go-project-pocs/src/handlers"
	"github.com/vfcoelho/go-project-pocs/src/repositories"
)

func SetCodeErrorMappings(mappings map[errs.ErrorCode]int) func(*fiber.Ctx) error { //REVIEW: fiber middleware to set error mappings and later be used by the error response middleware
	return func(c *fiber.Ctx) (err error) {
		c.Locals("codeErrorMappings", mappings)
		return c.Next()
	}
}

func ErrorRecoverMiddleware(c *fiber.Ctx) (err error) { //REVIEW: error response middleware to handle proper response - all errors will return a readable response to the caller
	err = c.Next()

	if err != nil {
		var httpErr *errs.Error
		var fiberErr *fiber.Error
		switch {
		case errors.As(err, &httpErr):
			errorCode, ok := c.Locals("codeErrorMappings").(map[errs.ErrorCode]int)[httpErr.Code]
			if !ok {
				errorCode = fiber.StatusInternalServerError
			}
			return c.Status(errorCode).JSON(httpErr)
		case errors.As(err, &fiberErr):
			return c.Status(fiberErr.Code).JSON(errs.NewError(err))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(errs.NewError(err))
		}
	}
	return err
}

func SetupRouter(app *fiber.App, producer handlers.EventProducer[dtos.Record]) {
	app.Use(recover.New())
	app.Use(ErrorRecoverMiddleware)
	app.Use(SetCodeErrorMappings(internal.MAPPING))

	memoryRepository := repositories.NewMemoryRepository[*dtos.Record]()

	app.Post("/v1/record", func(c *fiber.Ctx) error {
		return handlers.Post(c, memoryRepository, producer)
	})

	app.Get("/v1/record/:id", func(c *fiber.Ctx) error {
		return handlers.Get(c, memoryRepository)
	})
}
