package main

import (
	"log/slog"
	"os"
	"wb_tech/l3_1/internal/storage/postgres"
	"wb_tech/l3_1/internal/storage/queue"

	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	slog.Info("Starting worker")

	db, err := postgres.Connect()
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		return
	}

	repo := postgres.NewRepository(db)

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		slog.Error("RABBITMQ_URL is not set")
		return
	}

	consumer, err := queue.NewConsumer(rabbitURL, repo)
	if err != nil {
		slog.Error("Failed to create consumer", "error", err)
		return
	}
	defer consumer.Close()

	consumer.StartConsuming()
}
