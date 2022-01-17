package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"context"

	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func Shutdown(ctx context.Context, app *App, mongoClient *mongo.Client) error {
	return app.Shutdown(ctx)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Not loading Env var from file\n%s", err)
	}

	mongoClient := persistence.MongoDBInit()

	app, err := NewApp(mongoClient)
	if err != nil {
		log.Fatalf("Failed to prepare application\n%s", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	app.Run()

	<-done
	log.Println("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown gracefully")
	}
	log.Println("Server shutting down")
}
