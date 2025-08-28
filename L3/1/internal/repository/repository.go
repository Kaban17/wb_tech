package repository

import "wb_tech/l3_1/pkg/types"

type NotificationRepo interface {
	GetNotification(id int) (*types.Notification, error)
	CreateNotification(notification *types.Notification) error
	UpdateNotification(notification *types.Notification) error
	DeleteNotification(id int) error
}
