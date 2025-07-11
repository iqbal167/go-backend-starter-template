package main

import (
	"database/sql"
	"fmt"
	"go-backend-starter-template/internal/app/rest"
	"go-backend-starter-template/internal/config"
	"go-backend-starter-template/internal/database"
	"go-backend-starter-template/internal/pkg/logger"
	"go-backend-starter-template/internal/pkg/server"
	"go-backend-starter-template/internal/provider/postgres"
	"log/slog"
	"os"
	"time"

	"github.com/rs/cors"
)

// main is the entry point of the application.
func main() {
	if err := runApp(); err != nil {
		// Print the error message and Exit the program with a non-zero status code to indicate failure.
		fmt.Fprintf(os.Stderr, "Could not start application: %v\n", err)
		os.Exit(1)
	}
}

func runApp() error {
	config, err := config.Load()
	if err != nil {
		return fmt.Errorf("Failed to load config: %v\n", err)
	}

	logger := logger.New()

	db, err := newDB(config)
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %v\n", err)
	}
	defer db.Close()

	logger.Info("Database connection established")

	rest := rest.New(config)
	routes := rest.Routes()

	// Create a new server instance
	srv, err := server.NewServerBuilder().
		WithPort(config.App.Port).
		WithReadTimeout(5 * time.Second).
		WithWriteTimeout(10 * time.Second).
		WithLogLevel(slog.LevelInfo).
		WithCORS(cors.Options{
			AllowedOrigins:   config.CORSConfig.AllowedOrigins,
			AllowedMethods:   config.CORSConfig.AllowedMethods,
			AllowedHeaders:   config.CORSConfig.AllowedHeaders,
			AllowCredentials: config.CORSConfig.AllowCredentials,
			MaxAge:           config.CORSConfig.MaxAge,
		}).
		WithHandler(routes).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build server: %v\n", err)
	}

	if err = srv.Run(); err != nil {
		return fmt.Errorf("failed to run server: %v\n", err)
	}

	return nil
}

func newDB(config *config.Config) (*sql.DB, error) {
	dsn := postgres.NewDSN(postgres.PostgresConfig{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		User:     config.DB.User,
		Password: config.DB.Password,
		Name:     config.DB.Name,
		SSLMode:  config.DB.SSLMode,
	})

	return database.New(&database.Config{
		Driver: config.DB.Driver,
		DSN:    dsn,
	})
}
