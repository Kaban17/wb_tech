package cache

import (
	"orders/models"
	"sync"
)

type OrderCache struct {
	mu     sync.RWMutex
	orders map[string]models.Order
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		orders: make(map[string]models.Order),
	}
}

// Add добавляет заказ в кэш
func (c *OrderCache) Add(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.ID] = order
}

// Get возвращает заказ по ID
func (c *OrderCache) Get(id string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[id]
	return order, exists
}

// Preload загружает данные из БД при старте
func (c *OrderCache) Preload(orders []models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		c.orders[order.ID] = order
	}
}

// GetAll возвращает все заказы (для тестов)
func (c *OrderCache) GetAll() map[string]models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.orders
}
