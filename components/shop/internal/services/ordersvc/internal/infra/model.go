package infra

import (
	"time"

	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc/internal/domain"
)

type orderModel struct {
	ID         uint   `gorm:"primaryKey"`
	ProductID  uint   `gorm:"not null"`
	Quantity   int    `gorm:"not null"`
	TotalPrice int64  `gorm:"not null"`
	Status     string `gorm:"not null;default:'pending'"`
	CreatedAt  time.Time
}

func (*orderModel) TableName() string {
	return "orders"
}

func (m *orderModel) toDomain() *domain.Order {
	return &domain.Order{
		ID:         m.ID,
		ProductID:  m.ProductID,
		Quantity:   m.Quantity,
		TotalPrice: m.TotalPrice,
		Status:     m.Status,
		CreatedAt:  m.CreatedAt,
	}
}

func toOrderModel(d *domain.Order) *orderModel {
	return &orderModel{
		ID:         d.ID,
		ProductID:  d.ProductID,
		Quantity:   d.Quantity,
		TotalPrice: d.TotalPrice,
		Status:     d.Status,
		CreatedAt:  d.CreatedAt,
	}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&orderModel{})
}
