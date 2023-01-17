package models

import "time"

type Item struct {
	ID          uint64    `json:"id"`
	NameItem    string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Available   bool      `json:"available"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
