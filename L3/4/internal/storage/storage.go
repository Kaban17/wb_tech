package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) Save(ctx context.Context, id string, file io.Reader) error {
	path := filepath.Join(s.basePath, id)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}

func (s *LocalStorage) GetImageByID(ctx context.Context, id string) (io.ReadCloser, error) {
	path := filepath.Join(s.basePath, id)
	return os.Open(path)
}

func (s *LocalStorage) DeleteImage(ctx context.Context, id string) error {
	path := filepath.Join(s.basePath, id)
	return os.Remove(path)
}
