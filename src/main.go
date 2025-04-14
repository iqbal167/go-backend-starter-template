package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	http.HandleFunc("/healthz", healthHandler)

	addr := fmt.Sprintf("0.0.0.0:%d", 8080)

	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	log.Printf("Server started at http://%s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
