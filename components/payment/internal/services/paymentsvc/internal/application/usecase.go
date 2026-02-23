package application

import (
	"context"
	"fmt"

	"github.com/jed-jeon/go-echo-example/components/api"
	"github.com/jed-jeon/go-echo-example/components/payment/internal/common"
	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc/internal/domain"
)

// CreatePaymentCommand — 결제 생성 요청
type CreatePaymentCommand struct {
	OrderID   uint
	ProductID uint
	Quantity  int
	Amount    int64
	Method    string
}

// UseCase — paymentsvc 비즈니스 로직 계약
type UseCase interface {
	CreatePayment(ctx context.Context, cmd CreatePaymentCommand) (*domain.Payment, error)
}

type useCase struct {
	repo       PaymentRepository
	productFinder api.ProductFinder
}

func NewUseCase(repo PaymentRepository, productFinder api.ProductFinder) UseCase {
	return &useCase{
		repo:       repo,
		productFinder: productFinder,
	}
}

func (u *useCase) CreatePayment(ctx context.Context, cmd CreatePaymentCommand) (*domain.Payment, error) {
	product, err := u.productFinder.FindProduct(ctx, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify product: %w", err)
	}
	expectedAmount := product.Price * int64(cmd.Quantity)
	if expectedAmount != cmd.Amount {
		return nil, fmt.Errorf("expected %d, got %d: %w", expectedAmount, cmd.Amount, common.ErrAmountMismatch)
	}

	payment := &domain.Payment{
		OrderID:   cmd.OrderID,
		ProductID: cmd.ProductID,
		Amount:    cmd.Amount,
		Method:    cmd.Method,
		Status:    "completed",
	}

	if err = u.repo.Create(ctx, payment); err != nil {
		return nil, err
	}
	return payment, nil
}
