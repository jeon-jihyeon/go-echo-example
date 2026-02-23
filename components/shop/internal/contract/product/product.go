package product

import (
	"context"
	"time"
)

// Product — 컴포넌트 내 서비스 간 공유 도메인 타입
type Product struct {
	ID        uint
	Name      string
	Price     int64
	Stock     int
	CreatedAt time.Time
}

// Contract — productsvc가 같은 컴포넌트 내 서비스에 노출하는 계약
type Contract interface {
	FindProduct(ctx context.Context, id uint) (*Product, error)
}
