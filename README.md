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

## 의존성 흐름

### 수직 의존성

진입점에서 도메인까지의 의존 방향. `internal/`이 각 경계에서 접근 범위를 강제 차단.

```
cmd/api, cmd/command
        ↓
  components/              인터페이스 정의
        ↓
  components/api/          컴포넌트 간 공개 계약 (Level 1)
        ↓
  {component}.Component    컴포넌트 조립
        ↓
  {component}/internal/    컴포넌트 전용 (Level 2 internal)
        ↓
  {svc}/service.go         서비스 조립 (factory)
        ↓
  {svc}/internal/          서비스 전용 (Level 3 internal)
        ↓
  controller → application ← infra
                    ↓
                 domain
```

### 수평 의존성

컴포넌트 간, 서비스 간 의존은 반드시 계약(interface/DTO)을 경유.

```
컴포넌트 간    payment.Component → api.ProductFinder → shop.Component

서비스 간      ordersvc → contract.Product → productsvc
```

| 구분 | 계약 위치 | 계약 형태 | 예시 |
|------|-----------|-----------|------|
| 컴포넌트 간 | `components/api/` | interface + 단순 DTO | `ProductFinder`, `ProductCreator` |
| 서비스 간 | `{component}/internal/contract/` | struct (도메인 타입) | `contract.Product` |

### 서비스 내부 레이어

각 서비스의 `internal/` 내 클린 아키텍처 구성. 외부 의존성은 경계 레이어(controller, infra)에만 존재.

```
controller → application ← infra
                  ↓
               domain
```

| 레이어 | 역할 | 의존 대상 | 외부 의존성 |
|--------|------|-----------|------------|
| domain | 엔티티 정의 | 없음 | 없음 (stdlib만) |
| application | UseCase, Repository 인터페이스 | domain | 없음 |
| controller | HTTP 핸들링, DTO 변환 | application | echo, validator |
| infra | DB 접근, ORM 모델 변환 | application, domain | gorm |

### Internal 격리 규칙

Go `internal/` 패키지는 부모 디렉토리 트리 외부에서 import 불가. 이 규칙을 3단계로 중첩하여 의존성 범위를 강제 제한.

| Level | 위치 | 접근 가능 범위 | 차단 대상 |
|-------|------|---------------|-----------|
| 1 | `components/api/` | 모든 컴포넌트 | 제한 없음 (공개 계약) |
| 2 | `{component}/internal/` | 해당 컴포넌트 내 서비스 | 다른 컴포넌트 |
| 3 | `{svc}/internal/` | 해당 서비스 | 같은 컴포넌트의 다른 서비스 |

Level 3의 내부 코드는 `service.go`(factory 패턴)를 통해서만 외부 노출.

## 실행

```bash
go run cmd/api/main.go           # 서버 시작
go run cmd/command/main.go seed  # 샘플 데이터
```
