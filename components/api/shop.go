package api

import "context"

type Product struct {
	ID    uint
	Name  string
	Price int64
}

type ProductFinder interface {
	FindProduct(ctx context.Context, id uint) (*Product, error)
}

type CreateProductCommand struct {
	Name  string
	Price int64
	Stock int
}

type ProductCreator interface {
	CreateProduct(ctx context.Context, cmd CreateProductCommand) error
}
