package main

import (
	"fmt"
	"log/slog"
	"os"
	"wb_tech/l3_3/pkg/repository"

	"github.com/joho/godotenv"
)

var (
	DB_HOST     string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
)

func main() {
	fmt.Println("Hello, world from L3/3!")
	godotenv.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db, err := repository.Connect(DB_HOST, DB_PORT, DB_USER, DB_NAME, DB_PASSWORD)
	logger.Info("Connecting to database", "host", DB_HOST, "port", DB_PORT, "user", DB_USER, "name", DB_NAME)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	repository.NewCommentRepository(db)
}
