package models

import (
	"time"

	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
)

type Order struct {
	ID            uint64                 `json:"id"`
	OrderNumber   uint64                 `json:"order_number"`
	UserID        uint64                 `json:"user_id"`
	PromotionID   uint64                 `json:"promotion_id"`
	Subtotal      float64                `json:"subtotal"`
	TotalDiscount float64                `json:"total_discount"`
	Total         float64                `json:"total"`
	Quantity      uint64                 `json:"quantity"`
	Status        string                 `json:"status"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	Promotion     dto.PromotionDTOCreate `json:"promotion"`
	UserOrder     dto.UserDTOOrder       `json:"user"`
	Details       []OrderItems           `json:"detail_order"`
}
