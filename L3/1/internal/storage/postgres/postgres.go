package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	var err error
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
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
