package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manikandareas/go-clean-architecture/internal/delivery/http"
)

type RouteConfig struct {
	App                    *fiber.App
	BookController         *http.BookController
	UserController         *http.UserController
	AuthMiddleware         fiber.Handler
	RefreshTokenMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	api := c.App.Group("api")

	api.Post("/users", c.UserController.Register)
	api.Post("/users/_login", c.UserController.Login)
	api.Post("/users/_refresh", c.RefreshTokenMiddleware, c.UserController.RefreshToken)
}

func (c *RouteConfig) SetupAuthRoute() {
	api := c.App.Group("api", c.AuthMiddleware)
	api.Get("/books", c.BookController.FindAll)
	api.Post("/books", c.BookController.Create)
}
