package productsvc

import (
	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components/api"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/contract/product"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/application"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/controller"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc/internal/infra"
	"github.com/jed-jeon/go-echo-example/internal/core"
)

// Service — productsvc 서비스 인터페이스
type Service interface {
	NewController() core.Controller
	Migrate() error
	Contract() product.Contract
	Creator() api.ProductCreator
	Finder() api.ProductFinder
}

type service struct {
	db       *gorm.DB
	useCase  application.UseCase
	contract product.Contract
	creator  api.ProductCreator
	finder   api.ProductFinder
}

func NewService(db *gorm.DB) Service {
	repo := infra.NewProductRepository(db)
	uc := application.NewUseCase(repo)
	return &service{
		db:       db,
		useCase:  uc,
		contract: &contractAdapter{useCase: uc},
		creator:  &creator{useCase: uc},
		finder:   &finder{useCase: uc},
	}
}

func (s *service) NewController() core.Controller {
	return controller.NewController(s.useCase)
}

func (s *service) Migrate() error {
	return infra.Migrate(s.db)
}

func (s *service) Contract() product.Contract {
	return s.contract
}

func (s *service) Creator() api.ProductCreator {
	return s.creator
}

func (s *service) Finder() api.ProductFinder {
	return s.finder
}
