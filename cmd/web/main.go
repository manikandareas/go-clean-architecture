package main

import (
	"fmt"
	"github.com/manikandareas/go-clean-architecture/internal/config"
	"github.com/manikandareas/go-clean-architecture/pkg"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)

	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	jwtService := pkg.NewJwtService(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:         db,
		App:        app,
		Log:        log,
		Validate:   validate,
		Config:     viperConfig,
		JwtService: jwtService,
	})

	webPort := viperConfig.GetInt32("web.port")
	if err := app.Listen(fmt.Sprintf(":%d", webPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err.Error())
	}

}
