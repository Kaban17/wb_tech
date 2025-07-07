package service

import (
	"orders/internal/cache"
	"orders/internal/repository"
	"orders/models"
)

type OrderService struct {
	repo  *repository.OrderRepository
	cache *cache.OrderCache
}

func NewOrderService(repo *repository.OrderRepository, cache *cache.OrderCache) *OrderService {
	return &OrderService{repo: repo, cache: cache}
}

func (s *OrderService) GetOrderByUID(uid string) (*models.Order, error) {
	if order, exists := s.cache.Get(uid); exists {
		return &order, nil
	}

	order, err := s.repo.GetOrderByUID(uid)
	if err != nil {
		return nil, err
	}
	s.cache.Add(*order)
	return order, nil
}
func (s *OrderService) CreateOrder(order *models.Order) error {
	s.cache.Add(*order)
	return s.repo.CreateOrder(order)
}
