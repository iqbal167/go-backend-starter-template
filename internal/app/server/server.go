package server

import (
	"context"
	"database/sql"
	"fmt"
	"go-backend-starter-template/internal/config"
	"go-backend-starter-template/internal/database"
	"go-backend-starter-template/internal/provider/postgres"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	config   *config.Config
	database *sql.DB
}

func (s *Server) ServeHTTP(w *httptest.ResponseRecorder, req *http.Request) {
	panic("unimplemented")
}

func New(config *config.Config) *Server {
	return &Server{config: config}
}

func (s *Server) db() (*sql.DB, error) {
	dsn := postgres.NewDSN(postgres.PostgresConfig{
		Host:     s.config.DB.Host,
		Port:     s.config.DB.Port,
		User:     s.config.DB.User,
		Password: s.config.DB.Password,
		Name:     s.config.DB.Name,
		SSLMode:  s.config.DB.SSLMode,
	})

	db, err := database.New(&database.Config{
		Driver: s.config.DB.Driver,
		DSN:    dsn,
	})

	return db, err
}

func (s *Server) Run() {
	db, err := s.db()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	s.database = db

	addr := fmt.Sprintf("0.0.0.0:%d", 8080)
	router := s.routes()
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
