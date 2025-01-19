package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rfashwall/go-task/pkg/db"
	"github.com/rfashwall/go-task/pkg/middleware"
	"github.com/rfashwall/go-task/pkg/utils"
	"github.com/rfashwall/user-service/internal/command"
	"github.com/rfashwall/user-service/internal/handlers"
	"go.uber.org/zap"
)

func main() {
	shutdown := utils.InitTracer()
	defer shutdown()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(healthcheck.New())
	app.Use(middleware.TracingMiddleware("user-command-service"))

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// Connect to the database
	logger.Info("Connecting to MySQL")
	conn := db.MySqlConnect()
	defer conn.Close()

	logger.Info("Seeding data")
	db.SeedData(conn)

	// Initialize the repository
	userCommand := command.NewMySQLUserCommand(conn)

	// Initialize the handler
	logger.Debug("Initializing user command handler")
	userHandler := handlers.NewUserCommandHandler(userCommand, logger)

	// Set up routes
	logger.Info("Setting up routes")
	userHandler.SetupRoutes(app)

	logger.Info("Listening on port 3001")
	log.Fatal(app.Listen(":3001"))
}
