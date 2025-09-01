package main

import (
	"log/slog"
	"os"
	"wb_tech/l3_1/internal/api"
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

	slog.Info("Hello, world from L3/1!")
	db, err := postgres.Connect()
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		return
	}
	err = postgres.CreateTable(db)
	if err != nil {
		slog.Error("Error creating table", "error", err)
		return
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		slog.Error("RABBITMQ_URL is not set")
		return
	}

	s := api.NewServer(db, nil) // Producer will be set later

	// Create and start the consumer
	repo := postgres.NewRepository(db)
	consumer, err := queue.NewConsumer(rabbitURL, repo, s) // Pass server as broadcaster
	if err != nil {
		slog.Error("Failed to create consumer", "error", err)
		return
	}
	defer consumer.Close()
	go consumer.StartConsuming(0) // Start consuming in a separate goroutine

	producer, err := queue.NewProducer(rabbitURL)
	if err != nil {
		slog.Error("Failed to create producer", "error", err)
		return
	}
	defer producer.Close()

	s.SetProducer(producer) // Set the producer on the server

	r := api.NewRouter(s)
	slog.Info("starting server", "address", ":8080")
	api.Run(r, ":8080")
}
