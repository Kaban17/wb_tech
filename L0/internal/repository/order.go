package repository

import (
	"database/sql"
	"orders/internal/storage/database"
	"orders/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetOrderByUID(uid string) (*models.Order, error) {
	return database.GetOrder(r.db, uid)
}
func (r *OrderRepository) CreateOrder(order *models.Order) error {
	return database.CreateOrder(r.db, order)
}
