# Task 1 Brief: 统一 content-service 的分层结构

You are working in `C:\Users\Administrator\Desktop\personal web`.

Goal: reshape `services/content-service` to the approved layout:
`biz / cmd / idl / internal / pkg`, while keeping current behavior passing tests.

Ownership: only touch `services/content-service/**`, plus any docs/config files needed for that task. Do not edit auth-service or media-service.

Important requirements:
- `biz` is the HTTP entry layer.
- `internal/model` holds shared service-internal objects.
- `internal/service` holds business logic.
- `internal/dal/db` holds MySQL repositories.
- HTTP routing stays handwritten.
- Preserve article list/detail CRUD behavior.

Concrete work:
- Move HTTP handlers/routing/validation into `services/content-service/biz/article/*` and `services/content-service/biz/router.go`.
- Add `services/content-service/biz/model/article.go` for HTTP DTOs if needed.
- Add `services/content-service/internal/model/article.go` for service-internal objects.
- Move application logic into `services/content-service/internal/service/*`.
- Move repository code into `services/content-service/internal/dal/db/*`.
- Update `services/content-service/cmd/main.go` and `services/content-service/cmd/router.go` to wire the new structure.
- Keep existing tests green.

Validate with:
- `go test ./services/content-service/...`
- `go test ./...`
- if smoke is affected by your changes, keep `make micro-smoke` passing.

Report back:
- files changed
- tests run and result
- any concerns