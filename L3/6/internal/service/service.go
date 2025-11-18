package service

import (
	"context"
	"fmt"
	"time"

	"wb_tech/l3_6/internal/models"
	"wb_tech/l3_6/internal/repository"
)

// ItemService defines the interface for item-related business logic.
type ItemService interface {
	CreateItem(ctx context.Context, item *models.Item) error
	GetItem(ctx context.Context, id int) (*models.Item, error)
	UpdateItem(ctx context.Context, item *models.Item) error
	DeleteItem(ctx context.Context, id int) error
	ListItems(ctx context.Context, filters map[string]string) ([]models.Item, error)
	GetAnalytics(ctx context.Context, from, to time.Time) (*repository.AnalyticsResult, error)
}

// Service implements ItemService.
type Service struct {
	repo repository.ItemRepository
}

// NewItemService creates a new item service.
func NewItemService(repo repository.ItemRepository) ItemService {
	return &Service{repo: repo}
}

// CreateItem creates a new item after validating its data.
func (s *Service) CreateItem(ctx context.Context, item *models.Item) error {
	if err := validateItem(item); err != nil {
		return err
	}
	// Ensure ID is 0 for new items
	item.ID = 0
	return s.repo.CreateItem(ctx, item)
}

// GetItem retrieves an item by its ID.
func (s *Service) GetItem(ctx context.Context, id int) (*models.Item, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid item ID: %d", id)
	}
	return s.repo.GetItem(ctx, id)
}

// UpdateItem updates an existing item after validating its data.
func (s *Service) UpdateItem(ctx context.Context, item *models.Item) error {
	if item.ID <= 0 {
		return fmt.Errorf("invalid item ID for update: %d", item.ID)
	}
	if err := validateItem(item); err != nil {
		return err
	}
	return s.repo.UpdateItem(ctx, item)
}

// DeleteItem deletes an item by its ID.
func (s *Service) DeleteItem(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid item ID for deletion: %d", id)
	}
	return s.repo.DeleteItem(ctx, id)
}

// ListItems retrieves a list of items based on filters.
func (s *Service) ListItems(ctx context.Context, filters map[string]string) ([]models.Item, error) {
	return s.repo.ListItems(ctx, filters)
}

// GetAnalytics retrieves aggregated analytics for items within a date range.
func (s *Service) GetAnalytics(ctx context.Context, from, to time.Time) (*repository.AnalyticsResult, error) {
	if from.IsZero() || to.IsZero() || from.After(to) {
		return nil, fmt.Errorf("invalid date range: 'from' and 'to' dates are required and 'from' must be before or equal to 'to'")
	}
	return s.repo.GetAnalytics(ctx, from, to)
}

// validateItem performs basic validation on an Item.
func validateItem(item *models.Item) error {
	if item.Type == "" {
		return fmt.Errorf("item type cannot be empty")
	}
	if item.Amount < 0 {
		return fmt.Errorf("item amount cannot be negative")
	}
	if item.Date.IsZero() {
		return fmt.Errorf("item date cannot be empty")
	}
	if item.Category == "" {
		return fmt.Errorf("item category cannot be empty")
	}
	return nil
}
