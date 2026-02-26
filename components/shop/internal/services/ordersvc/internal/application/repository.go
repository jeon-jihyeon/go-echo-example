package application

import (
	"context"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc/internal/domain"
)

// OrderRepository — ordersvc 데이터 접근 계약
type OrderRepository interface {
	FindByID(ctx context.Context, id uint) (*domain.Order, error)
	Create(ctx context.Context, order *domain.Order) error
}
