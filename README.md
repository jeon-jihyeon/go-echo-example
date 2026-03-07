# 액터 모델 아키텍처

Actor 메시지 패싱으로 도메인 상태를 격리하는 아키텍처 예시.

## 핵심 개념

| 개념 | 설명 |
|------|------|
| 격리 단위 | Actor — mailbox가 메시지를 직렬 처리하여 동시성 안전 보장 |
| 통신 방식 | 타입 스위치 메시지 (`RequestFuture` → `Respond`) |
| 상태 관리 | Actor 내부 in-memory map — 외부 직접 접근 불가 |
| Actor 참조 | `*actor.PID` (untyped) — 컴파일 타임 메시지 타입 보장 없음 |

## 프로젝트 구조

```
cmd/api/main.go                  ActorSystem + Echo 부트스트랩

actors/
  registry.go                    Actor PID 일괄 생성/관리
  shop/                          상점 도메인
    messages.go                    메시지, 모델, 상수
    errors.go                      도메인 에러
    timeout.go                     Actor 타임아웃 상수
    product_actor.go               상품 CRUD
    order_actor.go                 주문 생성 (→ ProductActor)
  payment/                       결제 도메인
    messages.go                    메시지, 모델, 상수
    errors.go                      도메인 에러
    timeout.go                     Actor 타임아웃 상수
    payment_actor.go               결제 생성 (→ ProductActor)

handlers/                        HTTP → Actor 브릿지
  timeout.go                       핸들러 타임아웃 상수
  product_handler.go
  order_handler.go
  payment_handler.go

internal/core/
  response.go                    API 응답 포맷 + 유효성 검증
```

## Actor 계층

```
ActorSystem.Root
  ├─ ProductActor   ← 상품 상태 보유
  ├─ OrderActor     ← ProductActor PID 주입
  └─ PaymentActor   ← ProductActor PID 주입
```

## 통신 흐름

### HTTP → Actor

```
Client → Handler → root.RequestFuture(PID, msg, timeout)
                        → Actor.Receive(msg)
                        → ctx.Respond(result)
                   ← future.Result()
       ← JSON Response
```

### Actor → Actor (주문 생성)

```
POST /shop/orders
  → OrderHandler
    → root.RequestFuture(orderPID, CreateOrder)
      → OrderActor.Receive
        → ctx.RequestFuture(productPID, FindProduct)
          → ProductActor.Receive → ctx.Respond(ProductResult)
        ← 재고 확인 → 총액 계산 → 주문 저장
        → ctx.Respond(OrderResult)
    ← OrderResult
  → 201 Created
```

## API

| Method | Path | Actor | 설명 |
|--------|------|-------|------|
| POST | `/shop/products` | ProductActor | 상품 생성 |
| GET | `/shop/products` | ProductActor | 상품 목록 |
| GET | `/shop/products/:id` | ProductActor | 상품 조회 |
| POST | `/shop/orders` | OrderActor → ProductActor | 주문 생성 |
| GET | `/shop/orders/:id` | OrderActor | 주문 조회 |
| POST | `/payments` | PaymentActor → ProductActor | 결제 생성 |

### 응답 포맷

```json
{
  "code": 201,
  "message": "created",
  "data": {
    "id": 1,
    "product_id": 1,
    "quantity": 2,
    "total_price": 2000,
    "status": "pending",
    "created_at": "2026-03-07T12:00:00Z"
  }
}
```

## 알려진 트레이드오프

| 항목 | 설명 |
|------|------|
| `Receive` 내 블로킹 | `RequestFuture.Result()`가 Actor mailbox를 점유 — 프로덕션에서는 `PipeTo` 패턴 권장 |
| In-memory 상태 | Actor 재시작 시 상태 유실 — Event Sourcing 또는 Persistence 플러그인으로 해결 |
| Untyped PID | 컴파일 타임 메시지 타입 보장 없음 — 주석과 패키지 구조로 보완 |

## `internal/` 레이어드 아키텍처와 비교

| 항목 | internal 격리 | Actor 기반 |
|------|--------------|-----------|
| 격리 시점 | 컴파일 타임 (`internal/` 패키지) | 런타임 (Actor mailbox) |
| 서비스 간 통신 | 인터페이스 함수 호출 | `RequestFuture` 메시지 |
| 상태 저장 | DB + Repository 패턴 | Actor 내부 map |
| 동시성 제어 | goroutine + mutex | Actor mailbox 직렬화 |
| 의존성 방향 강제 | 컴파일러가 import 차단 | 메시지 타입으로 간접 참조 |

## 실행

```bash
go run cmd/api/main.go
```

```bash
# 상품 생성
curl -X POST localhost:8080/shop/products \
  -H 'Content-Type: application/json' \
  -d '{"name":"테스트","price":1000,"stock":10}'

# 주문 생성
curl -X POST localhost:8080/shop/orders \
  -H 'Content-Type: application/json' \
  -d '{"product_id":1,"quantity":2}'

# 결제
curl -X POST localhost:8080/payments \
  -H 'Content-Type: application/json' \
  -d '{"order_id":1,"product_id":1,"quantity":2,"amount":2000,"method":"card"}'
```

## 기술 스택

- [Proto.Actor](https://proto.actor/) — Go Actor 프레임워크
- [Echo v4](https://echo.labstack.com/) — HTTP 프레임워크
- [validator v10](https://github.com/go-playground/validator) — 구조체 유효성 검증
- In-memory 저장소 (DB 없음)
