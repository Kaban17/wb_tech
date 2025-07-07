package kafka

import (
	"context"
	"fmt"
	"log"
	"orders/internal/service"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type MessageHandler interface {
	Handle(message []byte) error
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  MessageHandler
}

func NewConsumer(brokers, groupID string, handler MessageHandler) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		handler:  handler,
	}, nil
}

func (c *Consumer) Subscribe(topics []string) error {
	return c.consumer.SubscribeTopics(topics, nil)
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				return err
			}

			if err := c.handler.Handle(msg.Value); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func InitKafkaConsumer(orderService *service.OrderService) (*Consumer, error) {
	brokers := os.Getenv("KAFKA_BROKERS")
	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "orders-service"
	}
	topic := os.Getenv("KAFKA_ORDERS_TOPIC")
	if topic == "" {
		topic = "orders"
	}

	// Создаем обработчик сообщений Kafka
	msgHandler := NewOrderMessageHandler(orderService)

	// Инициализируем потребителя
	consumer, err := NewConsumer(brokers, groupID, msgHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	// Подписываемся на топик
	if err := consumer.Subscribe([]string{topic}); err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	log.Printf("Kafka consumer initialized for topic %s", topic)
	return consumer, nil
}
