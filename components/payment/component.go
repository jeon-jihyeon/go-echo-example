package payment

import (
	"gorm.io/gorm"

	"github.com/jed-jeon/go-echo-example/components"
	"github.com/jed-jeon/go-echo-example/components/api"
	"github.com/jed-jeon/go-echo-example/components/payment/internal/services/paymentsvc"
)

// Component — payment 컴포넌트 인터페이스
type Component interface {
	components.Component
	SetProductFinder(finder api.ProductFinder)
}

type payment struct {
	db  *gorm.DB
	svc components.Service
}

func New(db *gorm.DB) Component {
	return &payment{db: db}
}

func (p *payment) Init() error {
	return nil
}

func (p *payment) Shutdown() error {
	return nil
}

func (p *payment) Migrate() error {
	return p.svc.Migrate()
}

func (p *payment) NewServer() components.RESTServer {
	return components.NewServer("", p.svc.NewController())
}

func (p *payment) SetProductFinder(finder api.ProductFinder) {
	p.svc = paymentsvc.NewService(p.db, finder)
}
