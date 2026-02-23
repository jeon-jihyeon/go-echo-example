package infra

import (
	"context"

	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc/internal/application"
	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc/internal/domain"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) application.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	m := toPaymentModel(payment)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	payment.ID = m.ID
	payment.CreatedAt = m.CreatedAt
	return nil
}
