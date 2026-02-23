package infra

import (
	"time"

	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc/internal/domain"
)

type paymentModel struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"not null"`
	OrderID   uint   `gorm:"not null"`
	Amount    int64  `gorm:"not null"`
	Method    string `gorm:"not null"`
	Status    string `gorm:"not null;default:'pending'"`
	CreatedAt time.Time
}

func (*paymentModel) TableName() string {
	return "payments"
}

func (m *paymentModel) toDomain() *domain.Payment {
	return &domain.Payment{
		ID:        m.ID,
		ProductID: m.ProductID,
		OrderID:   m.OrderID,
		Amount:    m.Amount,
		Method:    m.Method,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
	}
}

func toPaymentModel(d *domain.Payment) *paymentModel {
	return &paymentModel{
		ID:        d.ID,
		ProductID: d.ProductID,
		OrderID:   d.OrderID,
		Amount:    d.Amount,
		Method:    d.Method,
		Status:    d.Status,
		CreatedAt: d.CreatedAt,
	}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&paymentModel{})
}
