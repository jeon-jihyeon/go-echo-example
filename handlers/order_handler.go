package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"

	"github.com/jeon-jihyeon/go-echo-example/actors/shop"
	"github.com/jeon-jihyeon/go-echo-example/internal/core"
)

type OrderHandler struct {
	root     *actor.RootContext
	orderPID *actor.PID
}

func NewOrderHandler(root *actor.RootContext, orderPID *actor.PID) *OrderHandler {
	return &OrderHandler{root: root, orderPID: orderPID}
}

func (h *OrderHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/orders", h.create)
	g.GET("/orders/:id", h.getByID)
}

func (h *OrderHandler) create(c echo.Context) error {
	var req struct {
		ProductID uint `json:"product_id" validate:"required"`
		Quantity  int  `json:"quantity" validate:"required,gt=0"`
	}
	if err := core.BindAndValidate(c, &req); err != nil {
		return core.Error(c, http.StatusBadRequest, err.Error())
	}

	f := h.root.RequestFuture(h.orderPID, &shop.CreateOrder{
		ProductID: req.ProductID, Quantity: req.Quantity,
	}, actorTimeout)
	res, err := f.Result()
	if err != nil {
		return core.Error(c, http.StatusServiceUnavailable, "actor timeout")
	}

	result := res.(*shop.OrderResult)
	if result.Err != nil {
		switch {
		case errors.Is(result.Err, shop.ErrNotFound):
			return core.Error(c, http.StatusNotFound, "product not found")
		case errors.Is(result.Err, shop.ErrInsufficientStock):
			return core.Error(c, http.StatusBadRequest, "insufficient stock")
		default:
			return core.Error(c, http.StatusInternalServerError, result.Err.Error())
		}
	}
	return core.Created(c, toOrderResponse(result.Order))
}

func (h *OrderHandler) getByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, "invalid id")
	}

	f := h.root.RequestFuture(h.orderPID, &shop.FindOrder{ID: uint(id)}, actorTimeout)
	res, err := f.Result()
	if err != nil {
		return core.Error(c, http.StatusServiceUnavailable, "actor timeout")
	}
	result := res.(*shop.OrderResult)
	if result.Err != nil {
		return core.Error(c, http.StatusNotFound, "order not found")
	}
	return core.OK(c, toOrderResponse(result.Order))
}

type orderResponse struct {
	ID         uint      `json:"id"`
	ProductID  uint      `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice int64     `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func toOrderResponse(o *shop.Order) orderResponse {
	return orderResponse{
		ID: o.ID, ProductID: o.ProductID, Quantity: o.Quantity,
		TotalPrice: o.TotalPrice, Status: o.Status, CreatedAt: o.CreatedAt,
	}
}
