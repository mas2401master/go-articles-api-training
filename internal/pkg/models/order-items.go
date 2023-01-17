package models

import "time"

type OrderItems struct {
	ID        uint64    `json:"id"`
	OrderID   uint64    `json:"order_id"`
	ItemID    uint64    `json:"item_id"`
	NameItem  string    `json:"name_item"`
	Price     float64   `json:"price"`
	Quantity  uint64    `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
