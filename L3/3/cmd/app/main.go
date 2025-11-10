package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"wb_tech/l3_3/pkg/api/handlers"
	"wb_tech/l3_3/pkg/repository"
	"wb_tech/l3_3/pkg/service"

	"github.com/joho/godotenv"
)

var (
	POSTGRES_HOST     string
	POSTGRES_PORT     int
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	APP_PORT          int
)

func main() {
	fmt.Println("Hello, world from L3/3!")
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB = os.Getenv("POSTGRES_DB")
	POSTGRESPortStr := os.Getenv("POSTGRES_PORT")
	var err error
	POSTGRES_PORT, err = strconv.Atoi(POSTGRESPortStr)
	if err != nil {
		panic(fmt.Errorf("invalid DB_PORT: %w", err))
	}

	APP_PORTStr := os.Getenv("APP_PORT")
	APP_PORT, err = strconv.Atoi(APP_PORTStr)
	if err != nil {
		panic(fmt.Errorf("invalid APP_PORT: %w", err))
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db, err := repository.Connect(POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_DB, POSTGRES_PASSWORD)
	logger.Info("Connecting to database", "host", POSTGRES_HOST, "port", POSTGRES_PORT, "user", POSTGRES_USER, "db", POSTGRES_DB, "password", POSTGRES_PASSWORD)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()
	if err := db.CreateTable(); err != nil {
		panic(err)
	}
	repo := repository.NewCommentRepository(db)
	service := service.NewService(repo)
	api := handlers.NewAPI(service)
	logger.Info("Starting server on port", "port", APP_PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", APP_PORT), api.Router()); err != nil {
		panic(err)
	}
}
