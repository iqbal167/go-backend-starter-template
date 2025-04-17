package main

import (
	"context"
	"fmt"
	"go-backend-starter-template/src/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "app_name=%s version=%s", cfg.App.Name, cfg.App.Version)
	})

	addr := fmt.Sprintf("0.0.0.0:%d", 8080)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
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
