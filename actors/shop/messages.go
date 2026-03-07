package shop

import "time"

type ListProducts struct{}

type FindProduct struct {
	ID uint
}

type CreateProduct struct {
	Name  string
	Price int64
	Stock int
}

type Product struct {
	ID        uint
	Name      string
	Price     int64
	Stock     int
	CreatedAt time.Time
}

type ProductResult struct {
	Product *Product
	Err     error
}

type ProductListResult struct {
	Products []Product
	Err      error
}

type CreateOrder struct {
	ProductID uint
	Quantity  int
}

type FindOrder struct {
	ID uint
}

type Order struct {
	ID         uint
	ProductID  uint
	Quantity   int
	TotalPrice int64
	Status     string
	CreatedAt  time.Time
}

type OrderResult struct {
	Order *Order
	Err   error
}

const (
	OrderStatusPending = "pending"
)
