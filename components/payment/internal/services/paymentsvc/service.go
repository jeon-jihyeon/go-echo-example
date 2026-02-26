package paymentsvc

import (
	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components"
	"github.com/jed-jeon/go-echo-example/components/api"
	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc/internal/application"
	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc/internal/controller"
	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc/internal/infra"
	"github.com/jed-jeon/go-echo-example/internal/core"
)

type service struct {
	db      *gorm.DB
	useCase application.UseCase
}

func NewService(db *gorm.DB, productAPI api.ProductFinder) components.Service {
	repo := infra.NewPaymentRepository(db)
	uc := application.NewUseCase(repo, productAPI)
	return &service{db: db, useCase: uc}
}

func (s *service) NewController() core.Controller {
	return controller.NewController(s.useCase)
}

func (s *service) Migrate() error {
	return infra.Migrate(s.db)
}
