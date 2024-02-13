package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manikandareas/go-clean-architecture/internal/model"
	"github.com/manikandareas/go-clean-architecture/internal/usecase"
	"github.com/sirupsen/logrus"
)

type BookController struct {
	UseCase *usecase.BookUseCase
	Log     *logrus.Logger
}

func NewBookController(useCase *usecase.BookUseCase, log *logrus.Logger) *BookController {
	return &BookController{
		UseCase: useCase,
		Log:     log,
	}
}
func (c *BookController) FindAll(ctx *fiber.Ctx) error {
	response, err := c.UseCase.List(ctx.UserContext())

	if err != nil {
		c.Log.WithError(err).Error("failed to get books")
		return err
	}
	return ctx.JSON(fiber.Map{"data": response})
}

func (c *BookController) Create(ctx *fiber.Ctx) error {
	request := new(model.BookRequest)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create book")
		return err
	}
	return ctx.JSON(fiber.Map{"data": response})
}
