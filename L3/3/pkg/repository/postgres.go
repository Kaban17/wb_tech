package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS comments (
	id SERIAL PRIMARY KEY,
	parent_id INTEGER ,
	author_id INTEGER NOT NULL,
	text TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	path text NOT NULL,
	tsv TSVECTOR
);
DROP INDEX IF EXISTS idx_comments_tsv;
CREATE INDEX idx_comments_tsv ON comments USING GIN (tsv);
CREATE OR REPLACE FUNCTION comments_tsvector_trigger() RETURNS trigger AS $$
BEGIN
    -- 'russian' — конфигурация для русского языка (удаляет стоп-слова, выполняет стемминг)
    NEW.tsv := to_tsvector('english', NEW.text);
    RETURN NEW;
END
$$ LANGUAGE plpgsql;
DROP TRIGGER IF EXISTS tgr_comments_tsvector ON comments;
CREATE TRIGGER tgr_comments_tsvector
BEFORE INSERT OR UPDATE ON comments
FOR EACH ROW EXECUTE PROCEDURE comments_tsvector_trigger();
`

type DB struct {
	db *sqlx.DB
}

func Connect(host string, port int, user string, dbName string, password string) (*DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, password))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err) // ⬅️ Вернуть ошибку
	}
	return &DB{db}, nil
}
func (db *DB) CreateTable() error {
	db.db.MustExec(schema)
	return nil
}
func (db *DB) Close() error {
	return db.db.Close()
}
func (db *DB) Add(comment Comment, parentID *int) error {
	return db.db.QueryRowx(`INSERT INTO comments (parent_id, author_id, text, path) VALUES ($1, $2, $3, $4) RETURNING id`, parentID, comment.AuthorID, comment.Text, comment.Path).Scan(&comment.ID)
}
func (db *DB) Create(comment *Comment) error {
	return db.db.QueryRowx(`INSERT INTO comments (parent_id, author_id, text, path) VALUES ($1, $2, $3, $4) RETURNING id`, nil, comment.AuthorID, comment.Text, comment.Path).Scan(&comment.ID)
}
func (db *DB) GetThread(rootID *int, limit, offset int, sortField string, sortOrder string) ([]Comment, error) {
	var comments []Comment
	err := db.db.Select(&comments, `SELECT id, parent_id, author_id, text, created_at, path FROM comments WHERE parent_id = $1 ORDER BY `+sortField+` `+sortOrder+` LIMIT $2 OFFSET $3`, rootID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}
	return comments, nil
}
func (db *DB) Search(query string, limit, offset int) ([]Comment, error) {
	var comments []Comment
	err := db.db.Select(&comments, "SELECT id, parent_id, author_id, text, created_at, path FROM comments WHERE text ILIKE '%' || $1 || '%' ORDER BY created_at DESC LIMIT $2 OFFSET $3", query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search comments: %w", err)
	}
	return comments, nil
}
func (db *DB) DeleteComment(id string) error {
	_, err := db.db.Exec(`DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}
