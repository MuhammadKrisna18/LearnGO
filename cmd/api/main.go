package main

import (
	"log"

	"siakad-pro/config"
	"siakad-pro/internal/app"
)

func main() {
	log.Println("Starting application...")

	cfg := config.LoadConfig()

	application := app.NewApp(cfg)

	if err := application.Start(); err != nil {
		log.Fatalf("Fatal error starting application: %v\n", err)
	}
}
