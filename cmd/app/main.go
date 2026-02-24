package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/arunima10a/task-manager/config"
	"github.com/arunima10a/task-manager/internal/app"
)

// @title         Task Manager API
// @version       1.0
// @description   This is a clean architecture task management server.
// @host          localhost:8080
// @BasePath      /v1
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type 'Bearer <token>' to correctly authenticate
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")

	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
