package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/jed-jeon/go-echo-example/components/payment/internal/common"
	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc/internal/application"
	"github.com/jed-jeon/go-echo-example/internal/core"
)

type controller struct {
	core.ControllerBase
	useCase application.UseCase
}

func NewController(useCase application.UseCase) core.Controller {
	return &controller{useCase: useCase}
}

func (ctrl *controller) RegisterRoutes(router *echo.Group) {
	router.POST("/payments", ctrl.create)
}

func (ctrl *controller) create(c echo.Context) error {
	var req CreatePaymentRequest
	if err := ctrl.Bind(c, &req); err != nil {
		return core.Error(c, http.StatusBadRequest, err.Error())
	}

	payment, err := ctrl.useCase.CreatePayment(
		c.Request().Context(),
		application.CreatePaymentCommand{
			OrderID:   req.OrderID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Amount:    req.Amount,
			Method:    req.Method,
		},
	)
	if err != nil {
		if errors.Is(err, common.ErrAmountMismatch) {
			return core.Error(c, http.StatusBadRequest, "amount mismatch")
		}
		if errors.Is(err, common.ErrNotFound) {
			return core.Error(c, http.StatusNotFound, "product not found")
		}
		return core.Error(c, http.StatusInternalServerError, "internal server error")
	}
	return core.Created(c, toPaymentResponse(payment))
}
