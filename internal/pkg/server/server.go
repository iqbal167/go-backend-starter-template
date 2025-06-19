package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	shutdown   chan struct{}
}

func (s *Server) Run() error {
	serverErrors := make(chan error, 1)

	go func() {
		s.logger.Info("Server started", slog.String("address", fmt.Sprintf("http://%s", s.httpServer.Addr)))
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	// Wait for an interrupt signal OR a programmatic shutdown signal.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-osSignals:
		s.logger.Info("Shutdown signal received", slog.String("source", "OS Signal"), slog.String("signal", sig.String()))
		return s.gracefulShutdown()

	case <-s.shutdown:
		s.logger.Info("Shutdown signal received", slog.String("source", "Programmatic"))
		return s.gracefulShutdown()
	}
}

func (s *Server) gracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.httpServer.Close()
		return fmt.Errorf("could not stop server gracefully: %w", err)
	}

	s.logger.Info("Server stopped gracefully")
	return nil
}

// Shutdown initiates a programmatic shutdown of the server.
func (s *Server) Shutdown() {
	close(s.shutdown)
}
