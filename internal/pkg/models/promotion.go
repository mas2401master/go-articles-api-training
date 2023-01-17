package models

import "time"

type Promotion struct {
	ID        uint64    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Used      bool      `json:"used"`
	Discount  float64   `json:"discount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
