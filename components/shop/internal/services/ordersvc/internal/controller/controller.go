package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/common"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc/internal/application"
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
	router.POST("/orders", ctrl.create)
	router.GET("/orders/:id", ctrl.getByID)
}

func (ctrl *controller) create(c echo.Context) error {
	var req CreateOrderRequest
	if err := ctrl.Bind(c, &req); err != nil {
		return core.Error(c, http.StatusBadRequest, err.Error())
	}

	order, err := ctrl.useCase.CreateOrder(c.Request().Context(), req.ProductID, req.Quantity)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return core.Error(c, http.StatusNotFound, "product not found")
		}
		if errors.Is(err, common.ErrInsufficientStock) {
			return core.Error(c, http.StatusBadRequest, "insufficient stock")
		}
		return core.Error(c, http.StatusInternalServerError, "internal server error")
	}
	return core.Created(c, toOrderResponse(order))
}

func (ctrl *controller) getByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, "invalid id")
	}

	order, err := ctrl.useCase.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return core.Error(c, http.StatusNotFound, "order not found")
		}
		return core.Error(c, http.StatusInternalServerError, "internal server error")
	}
	return core.OK(c, toOrderResponse(order))
}
