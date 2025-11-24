package entity

import (
	"fmt"
	"time"
)

type format string
type status string
type Image struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	Format       format    `json:"format"`
	Status       status    `json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
	StoragePath  string    `json:"storage_path"`
}

func NewImage(originalName string, format format, status status) *Image {
	return &Image{
		ID:           GetUUID(),
		OriginalName: originalName,
		Format:       format,
		Status:       status,
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
		StoragePath:  fmt.Sprintf("%s_%s.%s", GetUUID(), string(status), format),
	}
}
