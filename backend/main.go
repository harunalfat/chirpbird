package main

import (
	"log"

	"github.com/harunalfat/chirpbird/backend/adapters/web"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load environment variables from file\n%s", err)
	}

	webServer, err := Initialize()
	if err != nil {
		log.Panicf("Failed to inject dependencies\n%s", err)
	}

	if err = web.Run(webServer); err != nil {
		log.Panicf("Failed to start web server\n%s", err)
	}
}
