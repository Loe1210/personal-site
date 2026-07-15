# Task 2 Brief: 统一 auth-service 的分层结构

You are working in `C:\Users\Administrator\Desktop\personal web`.

Goal: reshape `services/auth-service` to the approved layout:
`biz / cmd / idl / internal / pkg`, while keeping session-cookie-Redis login behavior intact.

Ownership: only touch `services/auth-service/**`, plus any docs/config files needed for that task. Do not edit content-service or media-service.

Important requirements:
- `biz` is the HTTP entry layer.
- `internal/model` holds shared service-internal objects.
- `internal/service` holds business logic.
- `internal/dal/db` holds MySQL repositories.
- HTTP routing stays handwritten.
- Login must continue using session + cookie + Redis.

Concrete work:
- Move HTTP handlers/routing/validation into `services/auth-service/biz/authenticator/*` and `services/auth-service/biz/router.go`.
- Add `services/auth-service/biz/model/auth.go` for HTTP DTOs if needed.
- Add `services/auth-service/internal/model/user.go` for service-internal objects.
- Move application logic into `services/auth-service/internal/service/*`.
- Move repository code into `services/auth-service/internal/dal/db/*`.
- Update `services/auth-service/cmd/main.go` and `services/auth-service/cmd/router.go` to wire the new structure.
- Keep existing tests green.

Validate with:
- `go test ./services/auth-service/...`
- `go test ./...`

Report back:
- files changed
- tests run and result
- any concerns