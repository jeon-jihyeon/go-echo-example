package domain

import "time"

type Product struct {
	ID        uint
	Name      string
	Price     int64
	Stock     int
	CreatedAt time.Time
}
