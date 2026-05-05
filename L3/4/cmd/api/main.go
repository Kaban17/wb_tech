package main
import (
	"database/sql"
	"log/slog"
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
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbConn := os.Getenv("DATABASE_URL")
	if dbConn == "" {
		dbConn = "postgres://user:password@localhost:5432/image_processor?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		slog.Error("Failed to connect to DB", "error", err)
		os.Exit(1)
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

	slog.Info("API server starting", "addr", ":8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

