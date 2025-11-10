package repository

import (
	"fmt"
	"time"
)

type Comment struct {
	ID        int       `json:"id" db:"id"`
	ParentID  *int      `json:"parent_id" db:"parent_id"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	Text      string    `json:"text" db:"text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Path      string    `json:"path" db:"path"`
}
type Repository struct {
	r DB
}

type CommentRepository interface {
	// Создание комментария, возвращает новый путь.
	Add(comment Comment, parentID *int) error

	// Получение комментария и всех его потомков (использует CTE).
	GetThread(rootID *int, limit, offset int, sortField string, sortOrder string) ([]Comment, error)

	// Удаление комментария и ВСЕХ потомков (использует оператор LIKE на path).
	DeleteSubtree(id int) error

	// Поиск по FTS (использует search_vector)
	Search(query string, limit, offset int) ([]Comment, error)
}

func NewCommentRepository(r *DB) CommentRepository {
	return &Repository{r: *r}
}

func (r *Repository) Add(comment Comment, parentID *int) error {
	return r.r.Create(&comment)
}

func (r *Repository) GetThread(rootID *int, limit, offset int, sortField string, sortOrder string) ([]Comment, error) {
	return r.r.GetThread(rootID, limit, offset, sortField, sortOrder)
}

func (r *Repository) DeleteSubtree(id int) error {
	return r.r.DeleteComment(fmt.Sprintf("%d%%", id))
}

func (r *Repository) Search(query string, limit, offset int) ([]Comment, error) {
	return r.r.Search(query, limit, offset)
}
