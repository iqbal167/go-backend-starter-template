package server

import (
	"context"
	"database/sql"
	"fmt"
	"go-backend-starter-template/internal/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	config *config.Config
	log    *slog.Logger
	db     *sql.DB
}

func New(config *config.Config, logger *slog.Logger, db *sql.DB) *Server {
	return &Server{config: config, log: logger, db: db}
}

func (s *Server) Run() {
	addr := fmt.Sprintf("0.0.0.0:%d", 8080)
	router := s.routes()
	handler := s.cors().Handler(router)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	serverErrors := make(chan error, 1)

	go func() {
		s.log.Info("Database connection established")
		s.log.Info("Server started", slog.String("address", fmt.Sprintf("http://%s", addr)))
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		s.log.Error("Server error", slog.String("error", err.Error()))

	case <-shutdown:
		s.log.Info("Starting shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			s.log.Error("Could not stop server gracefully", slog.String("error", err.Error()))
		}

		if err := s.db.Close(); err != nil {
			s.log.Error("Could not close database connection", slog.String("error", err.Error()))
		}

		s.log.Info("Server shut down gracefully")
	}
}
