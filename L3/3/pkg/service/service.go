package service

import (
	"wb_tech/l3_3/pkg/repository"
)

type CommentRepository interface {
	// Создание комментария, возвращает новый путь.
	Add(comment repository.Comment, parentID *int) error

	// Получение комментария и всех его потомков (использует CTE).
	GetThread(rootID *int, limit, offset int, sortField string, sortOrder string) ([]repository.Comment, error)

	// Удаление комментария и ВСЕХ потомков (использует оператор LIKE на path).
	DeleteSubtree(id int) error

	// Поиск по FTS (использует search_vector)
	Search(query string, limit, offset int) ([]repository.Comment, error)
}

type Service struct {
	repo CommentRepository
}

func NewService(repo CommentRepository) *Service {
	return &Service{repo: repo}
}
func (s *Service) Add(comment repository.Comment, parentID *int) error {
	return s.repo.Add(comment, parentID)
}
func (s *Service) Search(query string, limit, offset int) ([]repository.Comment, error) {
	return s.repo.Search(query, limit, offset)
}
func (s *Service) GetThread(rootID *int, limit, offset int, sortField string, sortOrder string) ([]repository.Comment, error) {
	return s.repo.GetThread(rootID, limit, offset, sortField, sortOrder)
}
func (s *Service) DeleteSubtree(id int) error {
	return s.repo.DeleteSubtree(id)
}
