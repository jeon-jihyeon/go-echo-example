package application

import (
	"context"

	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc/internal/domain"
)

// PaymentRepository — paymentsvc 데이터 접근 계약
type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
}
