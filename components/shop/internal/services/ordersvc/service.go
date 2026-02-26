package ordersvc

import (
	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/contract/product"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc/internal/application"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc/internal/controller"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc/internal/infra"
	"github.com/jed-jeon/go-echo-example/internal/core"
)

type service struct {
	db      *gorm.DB
	useCase application.UseCase
}

func NewService(db *gorm.DB, productContract product.Contract) components.Service {
	repo := infra.NewOrderRepository(db)
	uc := application.NewUseCase(repo, productContract)
	return &service{db: db, useCase: uc}
}

func (s *service) NewController() core.Controller {
	return controller.NewController(s.useCase)
}

func (s *service) Migrate() error {
	return infra.Migrate(s.db)
}
