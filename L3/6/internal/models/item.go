package models

import (
	"time"
)

type Item struct {
	ID          int       `json:"id" db:"id"`
	Type        string    `json:"type" db:"type"` // e.g., "income", "expense"
	Amount      float64   `json:"amount" db:"amount"`
	Date        time.Time `json:"date" db:"date"`
	Category    string    `json:"category" db:"category"`
	Description string    `json:"description" db:"description"`
}
