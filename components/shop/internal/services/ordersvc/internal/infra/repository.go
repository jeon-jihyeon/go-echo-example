package infra

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/common"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc/internal/application"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc/internal/domain"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) application.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindByID(ctx context.Context, id uint) (*domain.Order, error) {
	var m orderModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return m.toDomain(), nil
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	m := toOrderModel(order)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	order.ID = m.ID
	order.CreatedAt = m.CreatedAt
	return nil
}
