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
	*Http
	*Option
}

type Http struct {
	*http.Server
}

func (h *Http) Router(routes http.Handler) *Http {
	h.Server.Handler = routes
	return h
}

type Option struct {
	Port int
	Cors *Cors
	*slog.Logger
}

type Cors struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func New(routes http.Handler, opt *Option) *Server {
	// If no options are provided, create a default configuration.
	if opt == nil {
		opt = &Option{
			Port: 8080,
		}
	}

	// Initialize the logger if not provided.
	if opt.Logger == nil {
		opt.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	httpServer := newHttpServer(opt.Port)

	srv := &Server{
		Http:   httpServer,
		Option: opt,
	}

	// Apply CORS middleware if enabled.
	if opt.Cors != nil {
		routes = srv.cors().Handler(routes)
		srv.Info("CORS enabled", slog.Any("cors", opt.Cors))
	}

	// Set the final handler on the server.
	srv.Http.Router(routes)

	return srv
}

// newHttpServer creates an Http struct with an initialized server and logger.
func newHttpServer(port int) *Http {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	return &Http{
		Server: &http.Server{
			Addr:         addr,
			ErrorLog:     slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelError),
			WriteTimeout: 2 * time.Second,
		},
	}
}

func (srv *Server) Run() error {
	serverErrors := make(chan error, 1)

	go func() {
		srv.Info("Server started", slog.String("address", fmt.Sprintf("http://%s", srv.Http.Server.Addr)))
		serverErrors <- srv.Http.ListenAndServe()
	}()

	// Wait for an interrupt signal.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		srv.Info("Shutdown signal received", slog.String("signal", sig.String()))
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt a graceful shutdown.
		if err := srv.Http.Server.Shutdown(ctx); err != nil {
			srv.Http.Server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		srv.Info("Server stopped gracefully")
	}

	return nil
}
