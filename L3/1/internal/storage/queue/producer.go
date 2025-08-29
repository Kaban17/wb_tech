package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"wb_tech/l3_1/pkg/types"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewMQService(url string) (*MQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	err = setupDelayedExchange(ch)
	if err != nil {
		return nil, err
	}
	return &MQService{
		conn: conn,
		ch:   ch,
	}, nil
}

func setupDelayedExchange(ch *amqp.Channel) error {
	args := amqp.Table{
		"x-delayed-type": "direct",
	}

	err := ch.ExchangeDeclare(
		"notifications_delayed_exchange",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return fmt.Errorf("failed to declare delayed exchange: %w", err)
	}

	// Объявляем основную очередь для уведомлений
	_, err = ch.QueueDeclare(
		"notifications_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Привязываем очередь к exchange
	err = ch.QueueBind(
		"notifications_queue",
		"notification_routing_key",
		"notifications_delayed_exchange",
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}
func (r *MQService) Close() {
	r.ch.Close()
	r.conn.Close()
}
func (r *MQService) ScheduleNotification(notification *types.Notification) error {
	// Рассчитываем задержку до времени отправки
	now := time.Now()
	if notification.ScheduledAt.Before(now) {
		return fmt.Errorf("scheduled time is in the past")
	}

	delay := notification.ScheduledAt.Sub(now)
	delayMs := int64(delay / time.Millisecond)

	// Сериализуем уведомление в JSON
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Устанавливаем заголовок с задержкой
	headers := amqp.Table{
		"x-delay": delayMs,
	}

	// Публикуем сообщение с задержкой
	err = r.ch.PublishWithContext(
		context.Background(),
		"notifications_delayed_exchange",
		"notification_routing_key",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers:     headers,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}

	log.Printf("Notification scheduled for %s (in %v)",
		notification.ScheduledAt.Format(time.RFC3339), delay)

	return nil
}
