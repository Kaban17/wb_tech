package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"wb_tech/l3_4/internal/model"
	"wb_tech/l3_4/internal/processor"
	"wb_tech/l3_4/internal/repository"

	"github.com/segmentio/kafka-go"
	_ "github.com/lib/pq"
)

func main() {
	dbConn := os.Getenv("DATABASE_URL")
	if dbConn == "" {
		dbConn = "postgres://user:password@localhost:5432/image_processor?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	kafkaBrokers := []string{os.Getenv("KAFKA_BROKERS")}
	if kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}

	repo := repository.NewPostgresRepository(db)
	proc := processor.NewProcessor("./data")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   "image-tasks",
		GroupID: "image-processor-group",
	})
	defer reader.Close()

	log.Println("Worker started, waiting for tasks...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var task model.Task
		if err := json.Unmarshal(m.Value, &task); err != nil {
			log.Printf("Error unmarshaling task: %v", err)
			continue
		}

		log.Printf("Processing image: %s", task.ID)

		// Update status to processing
		if err := repo.UpdateStatusByID(context.Background(), task.ID, string(model.StatusProcessing)); err != nil {
			log.Printf("Error updating status: %v", err)
		}

		img := &model.Image{
			ID:           task.ID,
			OriginalPath: task.OriginalPath,
		}

		processedPath, err := proc.ProcessImage(context.Background(), img)
		if err != nil {
			log.Printf("Error processing image %s: %v", task.ID, err)
			repo.UpdateStatusByID(context.Background(), task.ID, string(model.StatusFailed))
			continue
		}

		if err := repo.UpdateProcessedPathByID(context.Background(), task.ID, processedPath); err != nil {
			log.Printf("Error updating processed path: %v", err)
		}

		log.Printf("Successfully processed image: %s", task.ID)
	}
}
