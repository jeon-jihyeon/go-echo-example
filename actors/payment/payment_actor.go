package payment

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/jeon-jihyeon/go-echo-example/actors/shop"
)

type PaymentActor struct {
	payments   map[uint]*Payment
	nextID     uint
	productPID *actor.PID
}

func NewPaymentActorProducer(productPID *actor.PID) func() actor.Actor {
	return func() actor.Actor {
		return &PaymentActor{
			payments:   make(map[uint]*Payment),
			productPID: productPID,
		}
	}
}

func (a *PaymentActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CreatePayment:
		f := ctx.RequestFuture(a.productPID, &shop.FindProduct{ID: msg.ProductID}, ActorTimeout)
		res, err := f.Result()
		if err != nil {
			ctx.Respond(&PaymentResult{Err: err})
			return
		}

		pr := res.(*shop.ProductResult)
		if pr.Err != nil {
			ctx.Respond(&PaymentResult{Err: ErrNotFound})
			return
		}

		expected := pr.Product.Price * int64(msg.Quantity)
		if expected != msg.Amount {
			ctx.Respond(&PaymentResult{Err: ErrAmountMismatch})
			return
		}

		a.nextID++
		p := &Payment{
			ID:        a.nextID,
			ProductID: msg.ProductID,
			OrderID:   msg.OrderID,
			Amount:    msg.Amount,
			Method:    msg.Method,
			Status:    PaymentStatusCompleted,
			CreatedAt: time.Now(),
		}
		a.payments[p.ID] = p
		ctx.Respond(&PaymentResult{Payment: p})
	}
}
