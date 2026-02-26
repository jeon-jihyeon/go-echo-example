package application

import (
	"context"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/domain"
)

// ProductRepository — productsvc 데이터 접근 계약
type ProductRepository interface {
	FindAll(ctx context.Context) ([]domain.Product, error)
	FindByID(ctx context.Context, id uint) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
}
