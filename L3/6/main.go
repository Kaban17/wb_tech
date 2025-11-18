package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"wb_tech/l3_6/internal/handlers"
	"wb_tech/l3_6/internal/repository"
	"wb_tech/l3_6/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Database connection details from environment variables
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPortStr := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	appPortStr := os.Getenv("APP_PORT")

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid POSTGRES_PORT: %v", err)
	}

	appPort, err := strconv.Atoi(appPortStr)
	if err != nil {
		log.Fatalf("Invalid APP_PORT: %v", err)
	}

	// Construct PostgreSQL connection string
	dbConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect to PostgreSQL
	db, err := sqlx.Connect("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	itemRepo := repository.NewPostgresRepository(db) // Corrected line

	// Initialize database schema
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := itemRepo.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Initialize service
	itemService := service.NewItemService(itemRepo)

	// Initialize handlers
	itemHandler := handlers.NewHandler(itemService)

	// Set up Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Register API routes
	itemHandler.RegisterRoutes(r)

	// Start HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", appPort),
		Handler: r,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %d", appPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %d: %v", appPort, err)
		}
	}()

	<-stop
	log.Println("Server shutting down...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}
