package entity

import "github.com/google/uuid"

const (
	StatusPending    status = "pending"
	StatusProcessing status = "processing"
	StatusCompleted  status = "completed"
	StatusFailed     status = "failed"
	PNG              format = "png"
	JPEG             format = "jpeg"
	JPG              format = "jpg"
	GIF              format = "gif"
)

func GetUUID() string {
	return uuid.New().String()
}
