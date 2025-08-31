package queue

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"wb_tech/l3_1/pkg/types"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName     = "notifications_exchange"
	queueName        = "notifications_queue"
	deadLetterQueue  = "notifications_wait"
	deadLetterExch   = "notifications_dlx"
	routingKey       = "notification_key"
)

type Producer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewProducer(url string) (*Producer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Exchange for final messages
	err = ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// Dead letter exchange for delayed messages
	err = ch.ExchangeDeclare(deadLetterExch, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// Queue for final messages
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// Waiting queue with DLX settings
	args := amqp.Table{
		"x-dead-letter-exchange":    exchangeName,
		"x-dead-letter-routing-key": routingKey,
	}
	_, err = ch.QueueDeclare(deadLetterQueue, true, false, false, false, args)
	if err != nil {
		return nil, err
	}

	// Bindings
	err = ch.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		return nil, err
	}
	err = ch.QueueBind(deadLetterQueue, routingKey, deadLetterExch, false, nil)
	if err != nil {
		return nil, err
	}

	return &Producer{conn: conn, ch: ch}, nil
}

func (p *Producer) Publish(notification *types.Notification, delay time.Duration) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.ch.PublishWithContext(ctx,
		deadLetterExch, // publish to dead-letter exchange
		routingKey,     // use the same routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Expiration:  strconv.Itoa(int(delay.Milliseconds())), // TTL for the message
		},
	)
}

func (p *Producer) Close() {
	p.ch.Close()
	p.conn.Close()
}