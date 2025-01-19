package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rfashwall/user-service/internal/handlers"
	"github.com/rfashwall/user-service/internal/repository"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// Database Configuration
	viper.SetEnvPrefix("USER_SERVICE") // Set environment variable prefix
	viper.AutomaticEnv()               // Enable Viper to read from environment variables.

	// Database Configuration from Environment Variables
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPass := viper.GetString("DB_PASS")
	dbName := viper.GetString("DB_NAME")

	// Construct the connection string
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Connect to the database
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Logger
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync() // flushes buffer, if any

	// Initialize Repository
	userRepo := repository.NewMySQLUserRepository(db)

	// Initialize Handler
	userHandler := handlers.NewUserHandler(userRepo, zapLogger)

	// Initialize Fiber App
	app := fiber.New()

	// Middleware
	app.Use(logger.New()) // Default Fiber logger

	// Routes
	app.Post("/users", userHandler.CreateUser)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Get("/users", userHandler.GetAllUsers)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
