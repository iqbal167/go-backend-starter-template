package server

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/rs/cors"
)

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	LogLevel     slog.Level
	CorsOptions  *cors.Options
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:         8080,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		LogLevel:     slog.LevelInfo,
		CorsOptions:  nil,
	}
}

func (c *ServerConfig) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("port must be a positive integer")
	}

	return nil
}

func BuildLogger(c *ServerConfig) *slog.Logger {
	handlerOptions := &slog.HandlerOptions{
		Level: c.LogLevel,
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
}
