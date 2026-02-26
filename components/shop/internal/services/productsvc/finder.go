package productsvc

import (
	"context"

	"github.com/jed-jeon/go-echo-example/components/api"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/application"
)

type finder struct {
	useCase application.UseCase
}

func (f *finder) FindProduct(ctx context.Context, id uint) (*api.Product, error) {
	p, err := f.useCase.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &api.Product{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
	}, nil
}
