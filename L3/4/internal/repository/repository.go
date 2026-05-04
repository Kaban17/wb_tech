package repository

import (
	"context"
	"database/sql"
	"wb_tech/l3_4/internal/model"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateImage(ctx context.Context, image *model.Image) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO images (id, original_path, processed_path, status) VALUES ($1, $2, $3, $4)",
		image.ID, image.OriginalPath, image.ProcessedPath, image.Status)
	return err
}

func (r *PostgresRepository) GetImageByID(ctx context.Context, id string) (*model.Image, error) {
	var img model.Image
	err := r.db.QueryRowContext(ctx,
		"SELECT id, original_path, processed_path, status, created_at FROM images WHERE id = $1", id).
		Scan(&img.ID, &img.OriginalPath, &img.ProcessedPath, &img.Status, &img.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *PostgresRepository) DeleteImage(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM images WHERE id = $1", id)
	return err
}

func (r *PostgresRepository) UpdateStatusByID(ctx context.Context, id string, status string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE images SET status = $1 WHERE id = $2", status, id)
	return err
}

func (r *PostgresRepository) UpdateProcessedPathByID(ctx context.Context, id string, path string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE images SET processed_path = $1, status = $2 WHERE id = $3", path, model.StatusCompleted, id)
	return err
}
