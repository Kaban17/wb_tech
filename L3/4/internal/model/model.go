package model

import "time"

type Image struct {
	ID            string    `json:"id"`
	OriginalPath  string    `json:"original_path"`
	ProcessedPath string    `json:"processed_path"`
	Status        Status    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
type Task struct {
	ID           string `json:"id"`
	OriginalPath string `json:"original_path"`
}
type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)
