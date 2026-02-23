package domain

import "time"

type Order struct {
	ID         uint
	ProductID  uint
	Quantity   int
	TotalPrice int64
	Status     string
	CreatedAt  time.Time
}
