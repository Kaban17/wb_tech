package kafka

import (
	"encoding/json"
	"log"
	"orders/internal/service"
	"orders/models"
)

type OrderMessageHandler struct {
	orderService *service.OrderService
}

func NewOrderMessageHandler(orderService *service.OrderService) *OrderMessageHandler {
	return &OrderMessageHandler{orderService: orderService}
}

func (h *OrderMessageHandler) Handle(message []byte) error {
	var order models.Order
	if err := json.Unmarshal(message, &order); err != nil {
		log.Printf("Failed to unmarshal order message: %v", err)
		return err
	}

	// Передаем заказ в сервис для обработки
	if err := h.orderService.CreateOrder(&order); err != nil {
		log.Printf("Failed to process order from Kafka: %v", err)
		return err
	}

	log.Printf("Successfully processed order from Kafka: %s", order.ID)
	return nil
}
