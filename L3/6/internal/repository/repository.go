package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver

	"wb_tech/l3_6/internal/models"
)

var schema = `
CREATE TABLE IF NOT EXISTS items (
	id SERIAL PRIMARY KEY,
	type VARCHAR(50) NOT NULL,
	amount NUMERIC(10, 2) NOT NULL CHECK (amount >= 0),
	date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	category VARCHAR(100) NOT NULL,
	description TEXT
);
`

// ItemRepository defines the interface for item data operations.
type ItemRepository interface {
	Init(ctx context.Context) error
	CreateItem(ctx context.Context, item *models.Item) error
	GetItem(ctx context.Context, id int) (*models.Item, error)
	UpdateItem(ctx context.Context, item *models.Item) error
	DeleteItem(ctx context.Context, id int) error
	ListItems(ctx context.Context, filters map[string]string) ([]models.Item, error)
	GetAnalytics(ctx context.Context, from, to time.Time) (*AnalyticsResult, error)
}

// postgresRepository implements ItemRepository using PostgreSQL.
type postgresRepository struct {
	db *sqlx.DB
}

// AnalyticsResult holds the aggregated analytics data.
type AnalyticsResult struct {
	Sum          float64 `json:"sum" db:"sum"`
	Average      float64 `json:"average" db:"average"`
	Count        int     `json:"count" db:"count"`
	Median       float64 `json:"median" db:"median"`
	Percentile90 float64 `json:"percentile90" db:"percentile_90"`
}

// NewPostgresRepository creates a new PostgreSQL item repository.
func NewPostgresRepository(db *sqlx.DB) ItemRepository {
	return &postgresRepository{db: db}
}

// Init initializes the database schema.
func (r *postgresRepository) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}
	return nil
}

// CreateItem inserts a new item into the database.
func (r *postgresRepository) CreateItem(ctx context.Context, item *models.Item) error {
	query := `INSERT INTO items (type, amount, date, category, description) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowxContext(ctx, query, item.Type, item.Amount, item.Date, item.Category, item.Description).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}
	return nil
}

// GetItem retrieves an item by its ID.
func (r *postgresRepository) GetItem(ctx context.Context, id int) (*models.Item, error) {
	var item models.Item
	query := `SELECT id, type, amount, date, category, description FROM items WHERE id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item with ID %d: %w", id, err)
	}
	return &item, nil
}

// UpdateItem updates an existing item in the database.
func (r *postgresRepository) UpdateItem(ctx context.Context, item *models.Item) error {
	query := `UPDATE items SET type = $1, amount = $2, date = $3, category = $4, description = $5 WHERE id = $6`
	result, err := r.db.ExecContext(ctx, query, item.Type, item.Amount, item.Date, item.Category, item.Description, item.ID)
	if err != nil {
		return fmt.Errorf("failed to update item with ID %d: %w", item.ID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected after updating item with ID %d: %w", item.ID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("item with ID %d not found for update", item.ID)
	}
	return nil
}

// DeleteItem deletes an item by its ID.
func (r *postgresRepository) DeleteItem(ctx context.Context, id int) error {
	query := `DELETE FROM items WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item with ID %d: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected after deleting item with ID %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("item with ID %d not found for deletion", id)
	}
	return nil
}

// ListItems retrieves a list of items based on filters.
func (r *postgresRepository) ListItems(ctx context.Context, filters map[string]string) ([]models.Item, error) {
	var items []models.Item
	query := `SELECT id, type, amount, date, category, description FROM items`

	whereClauses := []string{}
	args := []interface{}{}
	argCounter := 1

	if itemType, ok := filters["type"]; ok && itemType != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("type = $%d", argCounter))
		args = append(args, itemType)
		argCounter++
	}
	if category, ok := filters["category"]; ok && category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("category ILIKE $%d", argCounter))
		args = append(args, "%"+category+"%")
		argCounter++
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += ` ORDER BY date DESC`

	err := r.db.SelectContext(ctx, &items, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	return items, nil
}

// GetAnalytics calculates aggregated analytics for items within a given date range.
func (r *postgresRepository) GetAnalytics(ctx context.Context, from, to time.Time) (*AnalyticsResult, error) {
	var result AnalyticsResult
	query := `
		SELECT
			COALESCE(SUM(amount), 0) AS sum,
			COALESCE(AVG(amount), 0) AS average,
			COUNT(id) AS count,
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0) AS median,
			COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0) AS percentile_90
		FROM items
		WHERE date BETWEEN $1 AND $2
	`
	err := r.db.GetContext(ctx, &result, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics: %w", err)
	}
	return &result, nil
}
