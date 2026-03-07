package shop

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type ProductActor struct {
	products map[uint]*Product
	nextID   uint
}

func NewProductActor() actor.Actor {
	return &ProductActor{
		products: make(map[uint]*Product),
	}
}

func (a *ProductActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *ListProducts:
		list := make([]Product, 0, len(a.products))
		for _, p := range a.products {
			list = append(list, *p)
		}
		ctx.Respond(&ProductListResult{Products: list})

	case *FindProduct:
		p, ok := a.products[msg.ID]
		if !ok {
			ctx.Respond(&ProductResult{Err: ErrNotFound})
			return
		}
		ctx.Respond(&ProductResult{Product: p})

	case *CreateProduct:
		a.nextID++
		p := &Product{
			ID:        a.nextID,
			Name:      msg.Name,
			Price:     msg.Price,
			Stock:     msg.Stock,
			CreatedAt: time.Now(),
		}
		a.products[p.ID] = p
		ctx.Respond(&ProductResult{Product: p})
	}
}
