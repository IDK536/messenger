package app

import (
	"log"
	"messanger/internal/config"
	storage "messanger/internal/database"
	"messanger/internal/database/repositories"
	"messanger/internal/logger/sl"
	"messanger/internal/middleware"
	"messanger/internal/services"
	userController "messanger/internal/transport/http/user"

	"github.com/gofiber/fiber/v3"
)

func Run(cfg *config.Config) {
	logger := sl.InitLogger(cfg.Env)

	logger.Info("Logger is enabled")
	logger.Debug("Debug is enabled")

	database, err := storage.Connect(cfg, logger)

	if err != nil {
		log.Fatal(err)
	}

	defer database.Db.Close()

	app := fiber.New(fiber.Config{
		ReadTimeout: cfg.Timeout,
		IdleTimeout: cfg.IdleTimeout,
	})

	app.Use(middleware.NewLogger(logger))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello!")
	})

	repositories := repositories.InitRepositories(database)
	services := services.InitServices(repositories, cfg)

	userController := userController.InitUserController(services.UserService, services.TokenService, app)

	userController.Start(cfg)

	app.Listen(cfg.Address)
}
