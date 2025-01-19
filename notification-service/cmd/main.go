package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nats-io/nats.go"
	"github.com/rfashwall/go-task/pkg/middleware"
	"github.com/rfashwall/go-task/pkg/streaming"
	"github.com/rfashwall/go-task/pkg/utils"
	internalstream "github.com/rfashwall/notification-service/internal/streaming"
	nutils "github.com/rfashwall/notification-service/internal/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	shutdown := utils.InitTracer()
	defer shutdown()

	envPrefix := "NOTIFICATION_SERVICE"

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

	// Publish event to NATS
	natsService := viper.GetString("NATS_SERVICE")
	natsPort := viper.GetString("NATS_PORT")
	logger.Info("Connecting to NATS")
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", natsService, natsPort))
	if err != nil {
		log.Fatal("Failed to connect to NATS", err)
	}
	defer nc.Close()

	// Send email notification
	emailNotifier := nutils.NewEmailNotifier()
	eventHandler := internalstream.NewTaskEventHandler(emailNotifier, logger)

	// Initialize the NATS subscriber
	subscriber := streaming.NewNATSSubscriber(nc, eventHandler, logger)
	if err := subscriber.Subscribe("task.events"); err != nil {
		logger.Fatal("Failed to subscribe to NATS", zap.Error(err))
	}

	logger.Info(fmt.Sprintf("Listening on port %s", servicePort))
	log.Fatal(app.Listen(fmt.Sprintf(":%s", servicePort)))
}
