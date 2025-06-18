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

type Http struct {
	*http.Server
	*Option
	*slog.Logger
	shutdown chan struct{}
}

func (h *Http) Router(routes http.Handler) *Http {
	h.Server.Handler = routes
	return h
}

type Option struct {
	Port int
	Cors *Cors
}

type Cors struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func New(routes http.Handler, opt *Option) *Http {
	// If no options are provided, create a default configuration.
	if opt == nil {
		opt = &Option{
			Port: 8080,
		}
	}

	addr := fmt.Sprintf("0.0.0.0:%d", opt.Port)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	http := &Http{
		Server: &http.Server{
			Addr:     addr,
			ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		},
		Option:   opt,
		Logger:   logger,
		shutdown: make(chan struct{}),
	}

	// Apply CORS middleware if enabled.
	if opt.Cors != nil {
		routes = http.cors().Handler(routes)
		http.Info("CORS enabled", slog.Any("cors", opt.Cors))
	}

	// Set the final handler on the server.
	http.Router(routes)

	return http
}

func (http *Http) Run() error {
	serverErrors := make(chan error, 1)

	go func() {
		http.Info("Server started", slog.String("address", fmt.Sprintf("http://%s", http.Server.Addr)))
		serverErrors <- http.ListenAndServe()
	}()

	// Wait for an interrupt signal OR a programmatic shutdown signal.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-osSignals:
		http.Info("Shutdown signal received", slog.String("source", "OS Signal"), slog.String("signal", sig.String()))
		return http.gracefulShutdown()

	case <-http.shutdown:
		http.Info("Shutdown signal received", slog.String("source", "Programmatic"))
		return http.gracefulShutdown()
	}
}

func (http *Http) gracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := http.Server.Shutdown(ctx); err != nil {
		http.Server.Close()
		return fmt.Errorf("could not stop server gracefully: %w", err)
	}

	http.Info("Server stopped gracefully")
	return nil
}

// Shutdown initiates a programmatic shutdown of the server.
func (http *Http) Shutdown() {
	close(http.shutdown)
}
