package server

import (
	"context"
	"fmt"
	"go-backend-starter-template/internal/app/server/router"
	"go-backend-starter-template/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	config *config.Config
}

func New(config *config.Config) *Server {
	return &Server{config: config}
}

func (s *Server) Run() {
	router := router.New(s.config)

	addr := fmt.Sprintf("0.0.0.0:%d", 8080)
	handler := s.cors().Handler(router)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server started at http://%s", addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("server error: %v", err)

	case <-shutdown:
		log.Println("Starting shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			log.Fatalf("could not stop server gracefully: %v", err)
		}

		log.Printf("Server shut down gracefully")
	}
}
