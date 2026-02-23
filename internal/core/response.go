package core

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, APIResponse{
		Code:    http.StatusCreated,
		Message: "created",
		Data:    data,
	})
}

func Error(c echo.Context, status int, message string) error {
	return c.JSON(status, APIResponse{
		Code:    status,
		Message: message,
	})
}
