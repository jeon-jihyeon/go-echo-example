package application

import (
	"context"
	"fmt"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/common"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/contract/product"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc/internal/domain"
)

// UseCase — ordersvc 비즈니스 로직 계약
type UseCase interface {
	GetByID(ctx context.Context, id uint) (*domain.Order, error)
	CreateOrder(ctx context.Context, productID uint, quantity int) (*domain.Order, error)
}

type useCase struct {
	repo            OrderRepository
	productContract product.Contract
}

func NewUseCase(repo OrderRepository, productContract product.Contract) UseCase {
	return &useCase{
		repo:            repo,
		productContract: productContract,
	}
}

func (u *useCase) GetByID(ctx context.Context, id uint) (*domain.Order, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *useCase) CreateOrder(ctx context.Context, productID uint, quantity int) (*domain.Order, error) {
	p, err := u.productContract.FindProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product lookup failed: %w", err)
	}

	if p.Stock < quantity {
		return nil, common.ErrInsufficientStock
	}

	order := &domain.Order{
		ProductID:  productID,
		Quantity:   quantity,
		TotalPrice: p.Price * int64(quantity),
		Status:     "pending",
	}

	if err = u.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}
