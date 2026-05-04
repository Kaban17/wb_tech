package kafka

import (
	"context"
	"encoding/json"
	"wb_tech/l3_4/internal/model"

	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string, topic string) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Publisher) Publish(ctx context.Context, task *model.Task) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: payload,
	})
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
