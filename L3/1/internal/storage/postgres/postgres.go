package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = "5433"
	DB_USER     = "user"
	DB_PASSWORD = "user"
	DB_NAME     = "user"
)

func Connect() (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Close(db *sql.DB) error {
	return db.Close()
}
func CreateTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS Notifications (
		id SERIAL PRIMARY KEY,
		time_created TIMESTAMP DEFAULT NOW(),
		time_sent TIMESTAMP,
		scheduled_at TIMESTAMP,
		message TEXT NOT NULL,
		status TEXT,
		mail TEXT,
		tg TEXT

	);
	`
	_, err := db.Exec(query)
	if err != nil {
		// TODO log this error
	}
	fmt.Println("Table created successfully")
	return nil
}
