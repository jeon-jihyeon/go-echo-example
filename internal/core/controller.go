package core

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

// Controller — 모든 서비스 controller가 구현하는 라우트 등록 인터페이스
type Controller interface {
	RegisterRoutes(router *echo.Group)
}

// ControllerBase provides common Bind + struct validation helpers.
type ControllerBase struct{}

func (b *ControllerBase) Bind(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	return validate.Struct(req)
}
