package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/common"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/application"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/domain"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) application.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	var models []productModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	products := make([]domain.Product, len(models))
	for i, m := range models {
		products[i] = *m.toDomain()
	}
	return products, nil
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*domain.Product, error) {
	var m productModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return m.toDomain(), nil
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	m := toProductModel(product)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	product.ID = m.ID
	product.CreatedAt = m.CreatedAt
	return nil
}
