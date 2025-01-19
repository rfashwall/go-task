package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rfashwall/go-task/pkg/db"
	"github.com/rfashwall/go-task/pkg/middleware"
	"github.com/rfashwall/go-task/pkg/utils"
	"github.com/rfashwall/task-service/internal/handlers"
	"github.com/rfashwall/task-service/internal/repository/command"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	shutdown := utils.InitTracer()
	defer shutdown()

	envPrefix := "TASK_SERVICE"

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	serviceName := viper.GetString("SERVICE_NAME")
	servicePort := viper.GetString("SERVICE_PORT")

	app := fiber.New()
	app.Use(logger.New())
	app.Use(healthcheck.New())
	app.Use(middleware.TracingMiddleware(serviceName))

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// Connect to the database
	logger.Info("Connecting to MySQL")
	conn := db.MySqlConnect("TASK_SERVICE")
	defer conn.Close()

	// Initialize the repository
	taskCommand := command.NewMySQLTaskCommand(conn)

	// Initialize the handler
	logger.Info("Initializing task command handler")
	taskHandler := handlers.NewTaskCommandHandler(taskCommand, logger)

	// Set up routes
	logger.Info("Setting up routes")
	taskHandler.SetupRoutes(app)

	logger.Info(fmt.Sprintf("Listening on port %s", servicePort))
	log.Fatal(app.Listen(fmt.Sprintf(":%s", servicePort)))
}
