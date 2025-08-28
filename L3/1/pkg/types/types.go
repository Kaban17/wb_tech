package types

import (
	"time"
)

type Notification struct {
	TimeCreated time.Time          `json:"time_created"`
	Message     string             `json:"message"`
	Status      NotificationStatus `json:"status"`
	SendToMail  string             `json:"send_to_mail"`
	SentToTG    string             `json:"sent_to_tg"`
	TimeSent    time.Time          `json:"time_sent"`
	ScheduledAt time.Time          `json:"scheduled_at"`
}

type NotificationStatus string

const (
	Pending   NotificationStatus = "pending"
	Cancelled NotificationStatus = "cancelled"
	Delivered NotificationStatus = "delivered"
)

func NewNotification(message string, sendToMail string, sentToTG string, timeToSent time.Time) *Notification {
	return &Notification{
		TimeCreated: time.Now(),
		Message:     message,
		Status:      Pending,
		SendToMail:  sendToMail,
		SentToTG:    sentToTG,
		TimeSent:    timeToSent,
		ScheduledAt: timeToSent,
	}
}
