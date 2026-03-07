package actors

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/jeon-jihyeon/go-echo-example/actors/payment"
	"github.com/jeon-jihyeon/go-echo-example/actors/shop"
)

// Registry 모든 Actor의 PID를 보관
// *actor.PID는 Akka의 ActorRef<T>와 달리 untyped
//   - 컴파일 타임 메시지 타입 검증 불가
//   - 호출 측에서 각 PID에 올바른 메시지를 보내는 책임
//
// 도메인 패키지(shop/, payment/)의 메시지 타입으로 계약을 명시
type Registry struct {
	System     *actor.ActorSystem
	ProductPID *actor.PID // shop.ListProducts, shop.FindProduct, shop.CreateProduct
	OrderPID   *actor.PID // shop.CreateOrder, shop.FindOrder
	PaymentPID *actor.PID // payment.CreatePayment
}

func NewRegistry(system *actor.ActorSystem) *Registry {
	productProps := actor.PropsFromProducer(shop.NewProductActor)
	productPID, err := system.Root.SpawnNamed(productProps, "product-actor")
	if err != nil {
		log.Fatalf("failed to spawn product-actor: %v", err)
	}

	orderProps := actor.PropsFromProducer(shop.NewOrderActorProducer(productPID))
	orderPID, err := system.Root.SpawnNamed(orderProps, "order-actor")
	if err != nil {
		log.Fatalf("failed to spawn order-actor: %v", err)
	}

	paymentProps := actor.PropsFromProducer(payment.NewPaymentActorProducer(productPID))
	paymentPID, err := system.Root.SpawnNamed(paymentProps, "payment-actor")
	if err != nil {
		log.Fatalf("failed to spawn payment-actor: %v", err)
	}

	return &Registry{
		System:     system,
		ProductPID: productPID,
		OrderPID:   orderPID,
		PaymentPID: paymentPID,
	}
}
