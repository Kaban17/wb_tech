package types

import (
	"time"
)

type Notification struct {
	ID          int                `json:"id"`
	TimeCreated time.Time          `json:"time_created"`
	Message     string             `json:"message"`
	Status      NotificationStatus `json:"status"`
	Mail        string             `json:"mail"`
	TG          string             `json:"tg"`
	TimeSent    time.Time          `json:"time_sent"`
	ScheduledAt time.Time          `json:"scheduled_at"`
}

type NotificationStatus string

const (
	Pending   NotificationStatus = "pending"
	Cancelled NotificationStatus = "cancelled"
	Delivered NotificationStatus = "delivered"
)

func NewNotification(message string, mail string, tg string, timeToSent time.Time) *Notification {
	return &Notification{
		TimeCreated: time.Now(),
		Message:     message,
		Status:      Pending,
		Mail:        mail,
		TG:          tg,
		TimeSent:    timeToSent,
		ScheduledAt: timeToSent,
	}
}
