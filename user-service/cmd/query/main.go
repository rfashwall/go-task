package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rfashwall/go-task/pkg/db"
	"github.com/rfashwall/go-task/pkg/middleware"
	"github.com/rfashwall/go-task/pkg/utils"
	"github.com/rfashwall/user-service/internal/handlers"
	"github.com/rfashwall/user-service/internal/query"
	"go.uber.org/zap"
)

func main() {
	shutdown := utils.InitTracer()
	defer shutdown()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(healthcheck.New())
	app.Use(middleware.TracingMiddleware("user-query-service"))

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// Connect to the database
	logger.Info("Connecting to MySQL")
	conn := db.MySqlConnect()
	defer conn.Close()

	// Initialize the repository
	userQuery := query.NewMySQLUserQuery(conn)

	// Initialize the handler
	logger.Debug("Initializing user query handler")
	userHandler := handlers.NewUserQueryHandler(userQuery, logger)

	// Set up routes
	logger.Info("Setting up routes")
	userHandler.SetupRoutes(app)

	logger.Info("Listening on port 3000")
	log.Fatal(app.Listen(":3000"))
}
