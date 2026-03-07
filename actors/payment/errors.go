package payment

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrAmountMismatch = errors.New("amount mismatch")
)
