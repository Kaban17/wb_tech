package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"wb_tech/l3_4/internal/handler"
	"wb_tech/l3_4/internal/kafka"
	"wb_tech/l3_4/internal/repository"
	"wb_tech/l3_4/internal/service"
	"wb_tech/l3_4/internal/storage"

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
	store := storage.NewLocalStorage("./data")
	publisher := kafka.NewPublisher(kafkaBrokers, "image-tasks")
	defer publisher.Close()

	// Processor is nil for API side as it doesn't process images
	svc := service.NewService(store, repo, publisher, nil)
	h := handler.NewHandler(svc)

	mux := http.NewServeMux()
	handler.NewRouter(mux, h, "./static")

	fmt.Println("API server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
