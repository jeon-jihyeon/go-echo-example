package payment

import "time"

type CreatePayment struct {
	OrderID   uint
	ProductID uint
	Quantity  int
	Amount    int64
	Method    string
}

type Payment struct {
	ID        uint
	ProductID uint
	OrderID   uint
	Amount    int64
	Method    string
	Status    string
	CreatedAt time.Time
}

type PaymentResult struct {
	Payment *Payment
	Err     error
}

const (
	PaymentStatusCompleted = "completed"
)
