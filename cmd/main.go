package main

import (
	"fmt"
	"go-backend-starter-template/internal/app/server"
	"go-backend-starter-template/internal/config"
	"go-backend-starter-template/internal/database"
	"go-backend-starter-template/internal/pkg/logger"
	"go-backend-starter-template/internal/provider/postgres"
	"os"
)

// main is the entry point of the application.
func main() {
	config, err := config.Load()
	if err != nil {
		// Print the error message and Exit the program with a non-zero status code to indicate failure.
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	dsn := postgres.NewDSN(postgres.PostgresConfig{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		User:     config.DB.User,
		Password: config.DB.Password,
		Name:     config.DB.Name,
		SSLMode:  config.DB.SSLMode,
	})

	db, err := database.New(&database.Config{
		Driver: config.DB.Driver,
		DSN:    dsn,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	logger := logger.New()

	// Create a new server instance and start it.
	server := server.New(config, logger, db)
	server.Run()
}
