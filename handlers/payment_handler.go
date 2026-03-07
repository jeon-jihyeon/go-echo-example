package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"

	"github.com/jeon-jihyeon/go-echo-example/actors/payment"
	"github.com/jeon-jihyeon/go-echo-example/internal/core"
)

type PaymentHandler struct {
	root       *actor.RootContext
	paymentPID *actor.PID
}

func NewPaymentHandler(root *actor.RootContext, paymentPID *actor.PID) *PaymentHandler {
	return &PaymentHandler{root: root, paymentPID: paymentPID}
}

func (h *PaymentHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/payments", h.create)
}

func (h *PaymentHandler) create(c echo.Context) error {
	var req struct {
		OrderID   uint   `json:"order_id" validate:"required"`
		ProductID uint   `json:"product_id" validate:"required"`
		Quantity  int    `json:"quantity" validate:"required,gt=0"`
		Amount    int64  `json:"amount" validate:"required,gt=0"`
		Method    string `json:"method" validate:"required"`
	}
	if err := core.BindAndValidate(c, &req); err != nil {
		return core.Error(c, http.StatusBadRequest, err.Error())
	}

	f := h.root.RequestFuture(h.paymentPID, &payment.CreatePayment{
		OrderID: req.OrderID, ProductID: req.ProductID,
		Quantity: req.Quantity, Amount: req.Amount, Method: req.Method,
	}, actorTimeout)
	res, err := f.Result()
	if err != nil {
		return core.Error(c, http.StatusServiceUnavailable, "actor timeout")
	}

	result := res.(*payment.PaymentResult)
	if result.Err != nil {
		switch {
		case errors.Is(result.Err, payment.ErrNotFound):
			return core.Error(c, http.StatusNotFound, "product not found")
		case errors.Is(result.Err, payment.ErrAmountMismatch):
			return core.Error(c, http.StatusBadRequest, "amount mismatch")
		default:
			return core.Error(c, http.StatusInternalServerError, result.Err.Error())
		}
	}
	return core.Created(c, toPaymentResponse(result.Payment))
}

type paymentResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	OrderID   uint      `json:"order_id"`
	Amount    int64     `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func toPaymentResponse(p *payment.Payment) paymentResponse {
	return paymentResponse{
		ID: p.ID, ProductID: p.ProductID, OrderID: p.OrderID,
		Amount: p.Amount, Method: p.Method, Status: p.Status, CreatedAt: p.CreatedAt,
	}
}
