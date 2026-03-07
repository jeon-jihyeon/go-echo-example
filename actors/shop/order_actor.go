package shop

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type OrderActor struct {
	orders     map[uint]*Order
	nextID     uint
	productPID *actor.PID
}

func NewOrderActorProducer(productPID *actor.PID) func() actor.Actor {
	return func() actor.Actor {
		return &OrderActor{
			orders:     make(map[uint]*Order),
			productPID: productPID,
		}
	}
}

func (a *OrderActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *FindOrder:
		o, ok := a.orders[msg.ID]
		if !ok {
			ctx.Respond(&OrderResult{Err: ErrNotFound})
			return
		}
		ctx.Respond(&OrderResult{Order: o})

	case *CreateOrder:
		f := ctx.RequestFuture(a.productPID, &FindProduct{ID: msg.ProductID}, ActorTimeout)
		res, err := f.Result()
		if err != nil {
			ctx.Respond(&OrderResult{Err: err})
			return
		}

		pr := res.(*ProductResult)
		if pr.Err != nil {
			ctx.Respond(&OrderResult{Err: pr.Err})
			return
		}
		if pr.Product.Stock < msg.Quantity {
			ctx.Respond(&OrderResult{Err: ErrInsufficientStock})
			return
		}

		a.nextID++
		order := &Order{
			ID:         a.nextID,
			ProductID:  msg.ProductID,
			Quantity:   msg.Quantity,
			TotalPrice: pr.Product.Price * int64(msg.Quantity),
			Status:     OrderStatusPending,
			CreatedAt:  time.Now(),
		}
		a.orders[order.ID] = order
		ctx.Respond(&OrderResult{Order: order})
	}
}
