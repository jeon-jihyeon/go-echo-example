package application

import (
	"context"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/domain"
)

// CreateProductCommand — 상품 생성 요청
type CreateProductCommand struct {
	Name  string
	Price int64
	Stock int
}

// UseCase — productsvc 비즈니스 로직 계약
type UseCase interface {
	List(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id uint) (*domain.Product, error)
	Create(ctx context.Context, cmd CreateProductCommand) (*domain.Product, error)
}

type useCase struct {
	repo ProductRepository
}

func NewUseCase(repo ProductRepository) UseCase {
	return &useCase{repo: repo}
}

func (u *useCase) List(ctx context.Context) ([]domain.Product, error) {
	return u.repo.FindAll(ctx)
}

func (u *useCase) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *useCase) Create(ctx context.Context, cmd CreateProductCommand) (*domain.Product, error) {
	product := &domain.Product{
		Name:  cmd.Name,
		Price: cmd.Price,
		Stock: cmd.Stock,
	}
	if err := u.repo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}
