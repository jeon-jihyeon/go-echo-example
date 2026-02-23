package components

import (
	"github.com/labstack/echo/v4"

	"github.com/jed-jeon/go-echo-example/internal/core"
)

// Component — 모든 컴포넌트가 구현하는 lifecycle 인터페이스
type Component interface {
	Init() error
	Shutdown() error
	Migrate() error
	NewServer() RESTServer
}

// RESTServer — Echo 라우트 등록
type RESTServer interface {
	RoutePrefix() string
	RegisterRoutes(router *echo.Group)
}

// Service — 모든 서비스가 구현하는 공통 인터페이스
type Service interface {
	NewController() core.Controller
	Migrate() error
}

// server — RESTServer 공통 구현체
type server struct {
	routePrefix string
	controllers []core.Controller
}

func NewServer(routePrefix string, controllers ...core.Controller) RESTServer {
	return &server{routePrefix: routePrefix, controllers: controllers}
}

func (s *server) RoutePrefix() string {
	return s.routePrefix
}

func (s *server) RegisterRoutes(router *echo.Group) {
	for _, ctrl := range s.controllers {
		ctrl.RegisterRoutes(router)
	}
}
