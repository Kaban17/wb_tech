package sqlite

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/storage"

	sqlite3 "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS url (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL,
		count INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS stats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url_id INTEGER NOT NULL,
		accessed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		user_agent TEXT NOT NULL,
		FOREIGN KEY (url_id) REFERENCES url(id)
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url (alias, url) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	id, err := stmt.Exec(alias, urlToSave)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id.LastInsertId()
}
func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
func (s *Storage) UpdateURLStats(alias string, userAgent string) error {
	const op = "storage.sqlite.UpdateURLStats"

	stmt, err := s.db.Prepare("INSERT INTO stats (url_id, user_agent) VALUES ((SELECT id FROM url WHERE alias = ?), ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(alias, userAgent)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetDailyStats() ([]storage.Stat, error) {
	const op = "storage.sqlite.GetDailyStats"
	rows, err := s.db.Query(`
		SELECT u.url, u.alias, COUNT(s.id)
		FROM url u
		JOIN stats s ON u.id = s.url_id
		WHERE DATE(s.accessed_at) = DATE('now')
		GROUP BY u.id
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var stats []storage.Stat
	for rows.Next() {
		var stat storage.Stat
		if err := rows.Scan(&stat.URL, &stat.Alias, &stat.Count); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func (s *Storage) GetWeeklyStats() ([]storage.Stat, error) {
	const op = "storage.sqlite.GetWeeklyStats"
	rows, err := s.db.Query(`
		SELECT u.url, u.alias, COUNT(s.id)
		FROM url u
		JOIN stats s ON u.id = s.url_id
		WHERE DATE(s.accessed_at) >= DATE('now', '-7 days')
		GROUP BY u.id
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var stats []storage.Stat
	for rows.Next() {
		var stat storage.Stat
		if err := rows.Scan(&stat.URL, &stat.Alias, &stat.Count); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func (s *Storage) GetMonthlyStats() ([]storage.Stat, error) {
	const op = "storage.sqlite.GetMonthlyStats"
	rows, err := s.db.Query(`
		SELECT u.url, u.alias, COUNT(s.id)
		FROM url u
		JOIN stats s ON u.id = s.url_id
		WHERE DATE(s.accessed_at) >= DATE('now', '-1 month')
		GROUP BY u.id
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var stats []storage.Stat
	for rows.Next() {
		var stat storage.Stat
		if err := rows.Scan(&stat.URL, &stat.Alias, &stat.Count); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		stats = append(stats, stat)
	}

	return stats, nil
}
