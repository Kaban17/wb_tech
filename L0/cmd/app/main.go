package main

import (
	"context"
	"database/sql"
	"log"
	"orders/internal/cache"
	"orders/internal/handler"
	"orders/internal/kafka"
	"orders/internal/repository"
	"orders/internal/service"
	"orders/internal/storage/database"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Инициализация подключения к БД
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("Error closing database:", err)
		}
	}(db)

	log.Println("Connected to database")

	// Инициализация репозитория и сервиса
	c := cache.NewOrderCache()
	data, err := database.GetAllOrders(db)
	if err != nil {
		log.Fatal("Error getting orders:", err)
	}
	c.Preload(data)
	repo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(repo, c)

	// Создание HTTP сервера
	r := gin.Default()
	handler.InitOrderRouter(r, orderService)
	r.Static("/static", "./static")
	// Инициализация Kafka Consumer (если настроены переменные окружения)
	if os.Getenv("KAFKA_BROKERS") != "" {
		kafkaConsumer, err := kafka.InitKafkaConsumer(orderService)
		if err != nil {
			log.Fatal("Failed to initialize Kafka consumer:", err)
		}
		defer func() {
			if err := kafkaConsumer.Close(); err != nil {
				log.Println("Error closing Kafka consumer:", err)
			}
		}()

		// Запуск Kafka Consumer в отдельной горутине
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			if err := kafkaConsumer.Run(ctx); err != nil {
				log.Printf("Kafka consumer stopped with error: %v", err)
			}
		}()

		// Graceful shutdown для Kafka Consumer
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan
			cancel()
		}()
	} else {
		log.Println("KAFKA_BROKERS not set, Kafka consumer disabled")
	}

	// Запуск HTTP сервера
	log.Println("Starting HTTP server on :8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
