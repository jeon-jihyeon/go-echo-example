package domain

import "time"

type Payment struct {
	ID        uint
	ProductID uint
	OrderID   uint
	Amount    int64
	Method    string
	Status    string
	CreatedAt time.Time
}
