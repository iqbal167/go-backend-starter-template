package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/cors"
)

type ServerBuilder struct {
	config      *ServerConfig
	handler     http.Handler
	middlewares []func(http.Handler) http.Handler
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		config: DefaultServerConfig(),
	}
}

func (b *ServerBuilder) WithPort(port int) *ServerBuilder {
	b.config.Port = port
	return b
}

func (b *ServerBuilder) WithReadTimeout(duration time.Duration) *ServerBuilder {
	b.config.ReadTimeout = duration
	return b
}

func (b *ServerBuilder) WithWriteTimeout(duration time.Duration) *ServerBuilder {
	b.config.WriteTimeout = duration
	return b
}

func (b *ServerBuilder) WithIdleTimeout(duration time.Duration) *ServerBuilder {
	b.config.IdleTimeout = duration
	return b
}

func (b *ServerBuilder) WithCORS(option cors.Options) *ServerBuilder {
	b.config.CorsOptions = &option
	return b
}

func (b *ServerBuilder) WithHandler(handler http.Handler) *ServerBuilder {
	b.handler = handler
	return b
}

func (b *ServerBuilder) WithLogLevel(level slog.Level) *ServerBuilder {
	b.config.LogLevel = level
	return b
}

func (b *ServerBuilder) WithMiddleware(middleware func(http.Handler) http.Handler) *ServerBuilder {
	b.middlewares = append(b.middlewares, middleware)
	return b
}

func (b *ServerBuilder) Build() (*Server, error) {
	if err := b.config.Validate(); err != nil {
		return nil, fmt.Errorf("Failed to validate config: %w", err)
	}

	logger := BuildLogger(b.config)
	addr := fmt.Sprintf("0.0.0.0:%d", b.config.Port)

	handler := b.handler
	if handler == nil {
		return nil, fmt.Errorf("server handler must be provided using WithHandler")
	}

	for i := len(b.middlewares) - 1; i >= 0; i-- {
		handler = b.middlewares[i](handler)
	}

	// Apply CORS middleware if enabled
	if b.config.CorsOptions != nil {
		handler = cors.New(*b.config.CorsOptions).Handler(handler)
		logger.Info("CORS enabled")
	}

	httpServer := &http.Server{
		Addr:         addr,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		ReadTimeout:  b.config.ReadTimeout,
		WriteTimeout: b.config.WriteTimeout,
		IdleTimeout:  b.config.IdleTimeout,
		Handler:      handler,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
		shutdown:   make(chan struct{}),
	}, nil
}
