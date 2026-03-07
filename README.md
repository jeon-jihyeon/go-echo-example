# go-echo-example

Go와 Echo를 활용해 하나의 이커머스 도메인을 브랜치별로 다른 아키텍처로 구현한 프로젝트입니다.

## 아키텍처 목록

| 아키텍처 | 브랜치 | 핵심 개념 |
|----------|--------|-----------|
| 컴포넌트 계약 아키텍처 | [`feat/architecture-example`](https://github.com/jeon-jihyeon/go-echo-example/tree/feat/architecture-example) | Go `internal/` 패키지 규칙으로 컴파일 타임 의존성 격리 |
| 액터 모델 아키텍처 | [`feat/actor-example`](https://github.com/jeon-jihyeon/go-echo-example/tree/feat/actor-example) | Actor 메시지 패싱으로 런타임 상태 격리 |
