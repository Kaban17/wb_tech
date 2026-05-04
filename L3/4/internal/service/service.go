package service

import (
	"context"
	"io"
	"wb_tech/l3_4/internal/model"
)

type ImageStorage interface {
	Save(ctx context.Context, id string, file io.Reader) error
	GetImageByID(ctx context.Context, id string) (io.ReadCloser, error)
	DeleteImage(ctx context.Context, id string) error
}
type ImageRepository interface {
	CreateImage(ctx context.Context, image *model.Image) error
	GetImageByID(ctx context.Context, id string) (*model.Image, error)
	DeleteImage(ctx context.Context, id string) error
	UpdateStatusByID(ctx context.Context, id string, status string) error
	UpdateProcessedPathByID(ctx context.Context, id string, path string) error
}
type TaskPublisher interface {
	Publish(ctx context.Context, task *model.Task) error
}
type ImageProcessor interface {
	ProcessImage(ctx context.Context, image *model.Image) error
}
type service struct {
	imageStorage   ImageStorage
	imageRepo      ImageRepository
	taskPublisher  TaskPublisher
	imageProcessor ImageProcessor
}

func NewService(imageStorage ImageStorage, imageRepo ImageRepository, taskPublisher TaskPublisher, imageProcessor ImageProcessor) *service {
	return &service{
		imageStorage:   imageStorage,
		imageRepo:      imageRepo,
		taskPublisher:  taskPublisher,
		imageProcessor: imageProcessor,
	}
}
func (s *service) SaveImage(ctx context.Context, id string, file io.Reader) error {
	originalPath := id // For local storage, we can just use the ID as the filename
	if err := s.imageStorage.Save(ctx, originalPath, file); err != nil {
		return err
	}
	image := &model.Image{
		ID:           id,
		OriginalPath: originalPath,
		Status:       model.StatusPending,
	}
	if err := s.imageRepo.CreateImage(ctx, image); err != nil {
		return err
	}
	if err := s.taskPublisher.Publish(ctx, &model.Task{
		ID:           id,
		OriginalPath: originalPath,
	}); err != nil {
		return err
	}
	return nil
}
func (s *service) GetImageByID(ctx context.Context, id string) (*model.Image, error) {
	return s.imageRepo.GetImageByID(ctx, id)
}
func (s *service) DeleteImage(ctx context.Context, id string) error {
	image, err := s.imageRepo.GetImageByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.imageStorage.DeleteImage(ctx, image.OriginalPath); err != nil {
		// Ignore error if file doesn't exist
	}
	if image.ProcessedPath != "" {
		if err := s.imageStorage.DeleteImage(ctx, image.ProcessedPath); err != nil {
			// Ignore error
		}
	}
	if err := s.imageRepo.DeleteImage(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateProcessedPath(ctx context.Context, id string, path string) error {
	return s.imageRepo.UpdateProcessedPathByID(ctx, id, path)
}
