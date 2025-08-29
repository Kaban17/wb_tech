package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"wb_tech/l3_1/pkg/types"
)

func (r *MQService) StartConsumer(handler func(notification *types.Notification) error) error {
	msgs, err := r.ch.Consume(
		"notifications_queue",
		"notifications_consumer",
		false, // auto-ack = false (ручное подтверждение)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			var notification types.Notification
			err := json.Unmarshal(msg.Body, &notification)
			if err != nil {
				log.Printf("Failed to unmarshal notification: %v", err)
				msg.Nack(false, false) // Отклоняем сообщение без повторной очереди
				continue
			}

			// Обновляем время фактической отправки
			notification.TimeSent = time.Now()

			// Обрабатываем уведомление
			err = handler(&notification)
			if err != nil {
				log.Printf("Failed to handle notification: %v", err)
				msg.Nack(false, true) // Повторно ставим в очередь
				continue
			}

			// Обновляем статус
			notification.Status = types.Delivered

			log.Printf("Notification delivered: %s", notification.Message)
			msg.Ack(false) // Подтверждаем обработку
		}
	}()

	log.Println("Notification consumer started")
	return nil
}
