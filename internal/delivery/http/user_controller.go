package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manikandareas/go-clean-architecture/internal/delivery/http/middleware"
	"github.com/manikandareas/go-clean-architecture/internal/model"
	"github.com/manikandareas/go-clean-architecture/internal/usecase"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UseCase *usecase.UserUseCase
	Log     *logrus.Logger
}

func NewUserController(useCase *usecase.UserUseCase, log *logrus.Logger) *UserController {
	return &UserController{UseCase: useCase, Log: log}
}

func (u *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)

	if err := ctx.BodyParser(request); err != nil {
		u.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := u.UseCase.Register(ctx.Context(), request)
	if err != nil {
		u.Log.WithError(err).Error("failed to register user")
		return err
	}

	return ctx.JSON(fiber.Map{"data": response})
}

func (u *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)

	if err := ctx.BodyParser(request); err != nil {
		u.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	response, err := u.UseCase.Login(ctx.Context(), request)
	if err != nil {
		u.Log.WithError(err).Error("failed to login user")
		return err
	}

	return ctx.JSON(fiber.Map{"data": response})
}

func (u *UserController) RefreshToken(ctx *fiber.Ctx) error {
	request := middleware.GetUser(ctx)

	response, err := u.UseCase.RefreshToken(ctx.Context(), request)
	if err != nil {
		u.Log.WithError(err).Error("failed to refresh token")
		return err
	}

	return ctx.JSON(fiber.Map{"data": response})
}
