package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/manikandareas/go-clean-architecture/internal/delivery/http"
	"github.com/manikandareas/go-clean-architecture/internal/delivery/http/middleware"
	"github.com/manikandareas/go-clean-architecture/internal/delivery/http/route"
	"github.com/manikandareas/go-clean-architecture/internal/entity"
	"github.com/manikandareas/go-clean-architecture/internal/repository"
	"github.com/manikandareas/go-clean-architecture/internal/usecase"
	"github.com/manikandareas/go-clean-architecture/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB         *gorm.DB
	App        *fiber.App
	Log        *logrus.Logger
	Validate   *validator.Validate
	Config     *viper.Viper
	JwtService *pkg.JwtService
}

func Bootstrap(config *BootstrapConfig) {
	Migrator(config.DB, &entity.Book{}, &entity.User{})

	// setup	repository
	bookRepository := repository.NewBookRepository(config.Log)
	userRepository := repository.NewUserRepository(config.Log)
	// setup use case
	bookUseCase := usecase.NewBookUseCase(config.DB, config.Log, config.Validate, bookRepository)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, config.JwtService)
	//	setup controller
	bookController := http.NewBookController(bookUseCase, config.Log)
	userController := http.NewUserController(userUseCase, config.Log)
	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)
	refreshTokenMiddleware := middleware.NewRefreshToken(userUseCase)
	// setup route
	routeConfig := route.RouteConfig{
		App:                    config.App,
		BookController:         bookController,
		UserController:         userController,
		AuthMiddleware:         authMiddleware,
		RefreshTokenMiddleware: refreshTokenMiddleware,
	}
	routeConfig.Setup()
}
