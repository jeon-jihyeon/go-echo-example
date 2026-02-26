package controller

import (
	"time"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc/internal/domain"
)

type CreateOrderRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

type OrderResponse struct {
	ID         uint      `json:"id"`
	ProductID  uint      `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice int64     `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func toOrderResponse(o *domain.Order) *OrderResponse {
	return &OrderResponse{
		ID:         o.ID,
		ProductID:  o.ProductID,
		Quantity:   o.Quantity,
		TotalPrice: o.TotalPrice,
		Status:     o.Status,
		CreatedAt:  o.CreatedAt,
	}
}
