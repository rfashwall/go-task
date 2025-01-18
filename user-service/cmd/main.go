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
	"go.uber.org/zap"
)

func main() {
	// Database Configuration
	dbHost := "localhost"
	dbPort := "3306"
	dbUser := "root"
	dbPass := "root"
	dbName := "user_tracker"

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
	zapLogger, _ := zap.NewProduction()
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
