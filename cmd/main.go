package main

import (
	"go-backend-starter-template/internal/app/server"
	"go-backend-starter-template/internal/config"
)

func main() {
	config := config.Load()
	server := server.New(config)
	server.Run()
}
