package controller

import (
	"time"

	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc/internal/domain"
)

type CreatePaymentRequest struct {
	OrderID   uint   `json:"order_id" validate:"required"`
	ProductID uint   `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
	Amount    int64  `json:"amount" validate:"required,gt=0"`
	Method    string `json:"method" validate:"required"`
}

type PaymentResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	OrderID   uint      `json:"order_id"`
	Amount    int64     `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func toPaymentResponse(p *domain.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:        p.ID,
		ProductID: p.ProductID,
		OrderID:   p.OrderID,
		Amount:    p.Amount,
		Method:    p.Method,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
	}
}
