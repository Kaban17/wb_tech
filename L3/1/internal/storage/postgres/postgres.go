package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}
func GetDB() *sql.DB {
	return db
}
func Close() error {
	return db.Close()
}
func CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS Notifications (
		id SERIAL PRIMARY KEY,
		time_created TIMESTAMPZ DEFAULT NOW(),
		time_sent TIMESTAMPZ,
		scheduled_at TIMESTAMPZ,
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
	return nil
}
