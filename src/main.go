package main

import (
	"fmt"
	"go-backend-starter-template/src/config"
	"log"
	"net/http"
)

func main() {
	config := config.Load()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "app_name=%s version=%s", config.App.Name, config.App.Version)
	})

	addr := fmt.Sprintf("0.0.0.0:%d", 8080)

	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	log.Printf("Server started at http://%s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
