package productsvc

import (
	"context"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/contract/product"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/application"
)

type contractAdapter struct {
	useCase application.UseCase
}

func (a *contractAdapter) FindProduct(ctx context.Context, id uint) (*product.Product, error) {
	p, err := a.useCase.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &product.Product{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
	}, nil
}
