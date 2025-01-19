package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func MySqlConnect(envPrefix string) *sql.DB {

	// Database Configuration
	viper.SetEnvPrefix(envPrefix) // Set environment variable prefix
	viper.AutomaticEnv()          // Enable Viper to read from environment variables.

	// Database Configuration from Environment Variables
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPass := viper.GetString("DB_PASS")
	dbName := viper.GetString("DB_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return db
}

func SeedData(conn *sql.DB) {
	// Check if the users table is empty
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check users table: %v\n", err)
	}

	if count == 0 {
		// Insert sample data
		_, err = conn.Exec(`
			INSERT INTO users (name, email, password) VALUES
			('Alice Johnson', 'alice.johnson@example.com', 'hashed_password_1'),
			('Bob Smith', 'bob.smith@example.com', 'hashed_password_2'),
			('Charlie Brown', 'charlie.brown@example.com', 'hashed_password_3'),
			('Diana Prince', 'diana.prince@example.com', 'hashed_password_4'),
			('Ethan Hunt', 'ethan.hunt@example.com', 'hashed_password_5');
		`)
		if err != nil {
			log.Fatalf("Failed to seed users table: %v\n", err)
		}

		log.Println("Seed data inserted successfully.")
	} else {
		log.Println("Users table already seeded.")
	}
}
