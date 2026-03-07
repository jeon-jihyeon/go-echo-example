package core

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func BindAndValidate(c echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	return validate.Struct(req)
}

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OK(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, APIResponse{Code: http.StatusOK, Message: "success", Data: data})
}

func Created(c echo.Context, data any) error {
	return c.JSON(http.StatusCreated, APIResponse{Code: http.StatusCreated, Message: "created", Data: data})
}

func Error(c echo.Context, status int, msg string) error {
	return c.JSON(status, APIResponse{Code: status, Message: msg})
}
