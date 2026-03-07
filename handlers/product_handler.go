package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/labstack/echo/v4"

	"github.com/jeon-jihyeon/go-echo-example/actors/shop"
	"github.com/jeon-jihyeon/go-echo-example/internal/core"
)

type ProductHandler struct {
	root       *actor.RootContext
	productPID *actor.PID
}

func NewProductHandler(root *actor.RootContext, productPID *actor.PID) *ProductHandler {
	return &ProductHandler{root: root, productPID: productPID}
}

func (h *ProductHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/products", h.list)
	g.GET("/products/:id", h.getByID)
	g.POST("/products", h.create)
}

func (h *ProductHandler) list(c echo.Context) error {
	f := h.root.RequestFuture(h.productPID, &shop.ListProducts{}, actorTimeout)
	res, err := f.Result()
	if err != nil {
		return core.Error(c, http.StatusServiceUnavailable, "actor timeout")
	}
	result := res.(*shop.ProductListResult)
	return core.OK(c, toProductListResponse(result.Products))
}

func (h *ProductHandler) getByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, "invalid id")
	}

	f := h.root.RequestFuture(h.productPID, &shop.FindProduct{ID: uint(id)}, actorTimeout)
	res, err := f.Result()
	if err != nil {
		return core.Error(c, http.StatusServiceUnavailable, "actor timeout")
	}
	result := res.(*shop.ProductResult)
	if result.Err != nil {
		return core.Error(c, http.StatusNotFound, "product not found")
	}
	return core.OK(c, toProductResponse(result.Product))
}

func (h *ProductHandler) create(c echo.Context) error {
	var req struct {
		Name  string `json:"name" validate:"required"`
		Price int64  `json:"price" validate:"required,gt=0"`
		Stock int    `json:"stock" validate:"required,gte=0"`
	}
	if err := core.BindAndValidate(c, &req); err != nil {
		return core.Error(c, http.StatusBadRequest, err.Error())
	}

	f := h.root.RequestFuture(h.productPID, &shop.CreateProduct{
		Name: req.Name, Price: req.Price, Stock: req.Stock,
	}, actorTimeout)
	res, err := f.Result()
	if err != nil {
		return core.Error(c, http.StatusServiceUnavailable, "actor timeout")
	}
	result := res.(*shop.ProductResult)
	return core.Created(c, toProductResponse(result.Product))
}

type productResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}

func toProductResponse(p *shop.Product) productResponse {
	return productResponse{
		ID: p.ID, Name: p.Name, Price: p.Price,
		Stock: p.Stock, CreatedAt: p.CreatedAt,
	}
}

func toProductListResponse(ps []shop.Product) []productResponse {
	out := make([]productResponse, len(ps))
	for i, p := range ps {
		out[i] = productResponse{
			ID: p.ID, Name: p.Name, Price: p.Price,
			Stock: p.Stock, CreatedAt: p.CreatedAt,
		}
	}
	return out
}
