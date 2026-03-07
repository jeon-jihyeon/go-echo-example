package shop

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)
