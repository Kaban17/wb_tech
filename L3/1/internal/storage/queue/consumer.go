package queue

import (
	"encoding/json"
	"log/slog"
	"time"
	"wb_tech/l3_1/internal/repository"
	"wb_tech/l3_1/pkg/types"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broadcaster interface {
	BroadcastNotification(notification types.Notification, id int)
}

type Consumer struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	repo        repository.NotificationRepo
	broadcaster Broadcaster
}

func NewConsumer(url string, repo repository.NotificationRepo, broadcaster Broadcaster) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{conn: conn, ch: ch, repo: repo, broadcaster: broadcaster}, nil
}

func (c *Consumer) StartConsuming(delay time.Duration) {
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

			var notification types.Notification
			if err := json.Unmarshal(d.Body, &notification); err != nil {
				slog.Error("Failed to unmarshal notification", "error", err)
				continue
			}

			// Update the notification status to Delivered
			notification.Status = types.Delivered
			notification.TimeSent = time.Now()

			if err := c.repo.UpdateNotification(&notification, notification.ID); err != nil {
				slog.Error("Failed to update notification in db", "error", err, "id", notification.ID)
				continue
			}
			slog.Info("Successfully processed and updated notification", "id", notification.ID)

			// Broadcast the update to any connected SSE clients
			c.broadcaster.BroadcastNotification(notification, notification.ID)
		}
	}()

	slog.Info("Worker started. Waiting for messages.")
	<-forever
}

func (c *Consumer) Close() {
	c.ch.Close()
	c.conn.Close()
}
