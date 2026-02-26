package productsvc

import (
	"context"

	"github.com/jed-jeon/go-echo-example/components/api"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/application"
)

type creator struct {
	useCase application.UseCase
}

func (c *creator) CreateProduct(ctx context.Context, cmd api.CreateProductCommand) error {
	_, err := c.useCase.Create(ctx, application.CreateProductCommand{
		Name:  cmd.Name,
		Price: cmd.Price,
		Stock: cmd.Stock,
	})
	return err
}
