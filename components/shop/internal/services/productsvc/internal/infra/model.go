package infra

import (
	"time"

	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/domain"
)

type productModel struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Price     int64  `gorm:"not null"`
	Stock     int    `gorm:"not null;default:0"`
	CreatedAt time.Time
}

func (*productModel) TableName() string {
	return "products"
}

func (m *productModel) toDomain() *domain.Product {
	return &domain.Product{
		ID:        m.ID,
		Name:      m.Name,
		Price:     m.Price,
		Stock:     m.Stock,
		CreatedAt: m.CreatedAt,
	}
}

func toProductModel(d *domain.Product) *productModel {
	return &productModel{
		ID:        d.ID,
		Name:      d.Name,
		Price:     d.Price,
		Stock:     d.Stock,
		CreatedAt: d.CreatedAt,
	}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&productModel{})
}
