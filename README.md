# go-echo-example

Go `internal` 패키지 규칙을 활용한 3레벨 컴포넌트 아키텍처 예시.

## 전체 구조

```
cmd/api, cmd/command
    ↓
shop.Component
    productsvc → internal/ (application, controller, domain, infra)
    ordersvc   → internal/

payment.Component
    paymentsvc → internal/
```

### 의존성 흐름

```
shop.Component → ProductFinder() → payment.Component      components/api/
productsvc     → Contract()      → ordersvc               internal/contract/
```

- **컴포넌트 간** (`components/api/`) — 단순 DTO 계약
- **컴포넌트 내** (`internal/contract/`) — 도메인 타입 계약

### 3-Level Internal

| Level | 위치 | 접근 범위 |
|-------|------|----------|
| 1 | `components/api/` | 모든 컴포넌트 |
| 2 | `{component}/internal/` | 해당 컴포넌트 |
| 3 | `{svc}/internal/` | 해당 서비스 |

Level 3는 `service.go` (factory 패턴)로만 외부 노출.

## 실행

```bash
go run cmd/api/main.go           # 서버 시작
go run cmd/command/main.go seed  # 샘플 데이터
```
