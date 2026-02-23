package shop

import (
	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components"
	"github.com/jed-jeon/go-echo-example/components/api"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/ordersvc"
	"github.com/jed-jeon/go-echo-example/components/shop/internal/services/productsvc"
)

// Component — shop 컴포넌트 인터페이스
type Component interface {
	components.Component
	ProductCreator() api.ProductCreator
	ProductFinder() api.ProductFinder
}

type shop struct {
	db         *gorm.DB
	productSvc productsvc.Service
	orderSvc   components.Service
}

func New(db *gorm.DB) Component {
	productSvc := productsvc.NewService(db)
	orderSvc := ordersvc.NewService(db, productSvc.Contract())
	return &shop{
		db:         db,
		productSvc: productSvc,
		orderSvc:   orderSvc,
	}
}

func (s *shop) Init() error {
	return nil
}

func (s *shop) Shutdown() error {
	return nil
}

func (s *shop) Migrate() error {
	if err := s.productSvc.Migrate(); err != nil {
		return err
	}
	return s.orderSvc.Migrate()
}

func (s *shop) NewServer() components.RESTServer {
	return components.NewServer("/shop",
		s.productSvc.NewController(),
		s.orderSvc.NewController(),
	)
}

func (s *shop) ProductCreator() api.ProductCreator {
	return s.productSvc.Creator()
}

func (s *shop) ProductFinder() api.ProductFinder {
	return s.productSvc.Finder()
}
