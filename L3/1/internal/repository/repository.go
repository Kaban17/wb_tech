package repository

import (
	"wb_tech/l3_1/pkg/types"
)

type NotificationRepo interface {
	GetNotification(id int) (*types.Notification, error)
	CreateNotification(notification *types.Notification) (int, error)
	UpdateNotification(notification *types.Notification, id int) error
	DeleteNotification(id int) error
}
