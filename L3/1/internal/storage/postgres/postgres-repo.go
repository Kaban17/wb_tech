package postgres

import (
	"database/sql"
	"time"
	"wb_tech/l3_1/pkg/types"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
func (r *Repository) CreateNotification(notification *types.Notification) (int, error) {
	query := `
	INSERT INTO Notifications (
		time_created,
		time_sent,
		scheduled_at,
		message,
		status,
		mail,
		tg
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	)
	RETURNING id
	`
	var id int
	err := r.db.QueryRow(query,
		time.Now(),
		notification.TimeSent,
		notification.ScheduledAt,
		notification.Message,
		types.Pending,
		notification.Mail,
		notification.TG,
	).Scan(&id)
	return id, err
}
func (r *Repository) GetNotification(id int) (*types.Notification, error) {
	query := `
	SELECT
		time_created,
		time_sent,
		scheduled_at,
		message,
		status,
		mail,
		tg
	FROM Notifications
	WHERE id = $1
	`
	var notification types.Notification
	err := r.db.QueryRow(query, id).Scan(
		&notification.TimeCreated,
		&notification.TimeSent,
		&notification.ScheduledAt,
		&notification.Message,
		&notification.Status,
		&notification.Mail,
		&notification.TG,
	)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}
func (r *Repository) UpdateNotification(notification *types.Notification, id int) error {
	query := `
	UPDATE Notifications
	SET
		time_created = $2,
		time_sent = $3,
		scheduled_at = $4,
		message = $5,
		status = $6,
		mail = $7,
		tg = $8
	WHERE id = $1
	`
	_, err := r.db.Exec(query,
		id,
		notification.TimeCreated,
		notification.TimeSent,
		notification.ScheduledAt,
		notification.Message,
		notification.Status,
		notification.Mail,
		notification.TG,
	)
	return err
}
func (r *Repository) DeleteNotification(id int) error {
	query := `
	DELETE FROM Notifications
	WHERE id = $1
	`
	_, err := r.db.Exec(query, id)
	return err
}
