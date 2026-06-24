package main

import (
	"log"

	"modular-monolith/config"
	"modular-monolith/internal/app"
)

// @title Modular Monolith API
// @version 1.0
// @description This is a sample API for a Modular Monolith built with Go Fiber.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log.Println("Starting application...")

	cfg := config.LoadConfig()

	application := app.NewApp(cfg)
	if err := application.Start(); err != nil {
		log.Fatalf("Fatal error starting application: %v\n", err)
	}
}
