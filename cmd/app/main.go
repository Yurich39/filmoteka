package main

import (
	"log"

	"people-finder/config"
	"people-finder/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	log.Println("Config:\n", cfg)

	// Run
	app.Run(cfg)
}
