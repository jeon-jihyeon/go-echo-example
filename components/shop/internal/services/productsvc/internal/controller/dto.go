package controller

import (
	"time"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/domain"
)

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required"`
	Price int64  `json:"price" validate:"required,gt=0"`
	Stock int    `json:"stock" validate:"required,gte=0"`
}

type ProductResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}

func toProductResponse(p *domain.Product) *ProductResponse {
	return &ProductResponse{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
	}
}

func toProductListResponse(products []domain.Product) []ProductResponse {
	res := make([]ProductResponse, len(products))
	for i := range products {
		res[i] = *toProductResponse(&products[i])
	}
	return res
}
