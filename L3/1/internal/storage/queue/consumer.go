package queue

import (
	"encoding/json"
	"log/slog"
	"time"
	"wb_tech/l3_1/internal/repository"
	"wb_tech/l3_1/pkg/types"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	repo repository.NotificationRepo
}

func NewConsumer(url string, repo repository.NotificationRepo) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{conn: conn, ch: ch, repo: repo}, nil
}

func (c *Consumer) StartConsuming() {
	msgs, err := c.ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		slog.Error("Failed to register a consumer", "error", err)
		return
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			slog.Info("Received a message", "body", string(d.Body))

			// Temporary struct to unmarshal the message from the queue
			var rawMsg struct {
				Message string `json:"message"`
				Mail    string `json:"mail"`
				TG      string `json:"tg"`
			}

			if err := json.Unmarshal(d.Body, &rawMsg); err != nil {
				slog.Error("Failed to unmarshal notification", "error", err)
				continue
			}

			// Use the constructor to create a valid notification object
			notification := types.NewNotification(rawMsg.Message, rawMsg.Mail, rawMsg.TG, time.Now()) // Assuming ScheduledAt is now

			if _, err := c.repo.CreateNotification(notification); err != nil {
				slog.Error("Failed to create notification in db", "error", err)
				continue
			}
			slog.Info("Successfully processed and saved notification")
		}
	}()

	slog.Info("Worker started. Waiting for messages.")
	<-forever
}

func (c *Consumer) Close() {
	c.ch.Close()
	c.conn.Close()
}
