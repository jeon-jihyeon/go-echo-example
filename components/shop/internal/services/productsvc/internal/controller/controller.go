package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/jed-jeon/go-echo-example/components/shop/internal/common"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/application"
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
	router.GET("/products", ctrl.list)
	router.GET("/products/:id", ctrl.getByID)
	router.POST("/products", ctrl.create)
}

func (ctrl *controller) list(c echo.Context) error {
	products, err := ctrl.useCase.List(c.Request().Context())
	if err != nil {
		return core.Error(c, http.StatusInternalServerError, "internal server error")
	}
	return core.OK(c, toProductListResponse(products))
}

func (ctrl *controller) getByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return core.Error(c, http.StatusBadRequest, "invalid id")
	}

	product, err := ctrl.useCase.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return core.Error(c, http.StatusNotFound, "product not found")
		}
		return core.Error(c, http.StatusInternalServerError, "internal server error")
	}
	return core.OK(c, toProductResponse(product))
}

func (ctrl *controller) create(c echo.Context) error {
	var req CreateProductRequest
	if err := ctrl.Bind(c, &req); err != nil {
		return core.Error(c, http.StatusBadRequest, err.Error())
	}

	product, err := ctrl.useCase.Create(c.Request().Context(), application.CreateProductCommand{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	})
	if err != nil {
		return core.Error(c, http.StatusInternalServerError, "internal server error")
	}
	return core.Created(c, toProductResponse(product))
}
