# Task 3 Brief: 统一 media-service 的分层结构

You are working in `C:\Users\Administrator\Desktop\personal web`.

Goal: reshape `services/media-service` to the approved layout:
`biz / cmd / idl / internal / pkg`, while keeping upload and file lookup behavior intact.

Ownership: only touch `services/media-service/**`, plus any docs/config files needed for that task. Do not edit auth-service or content-service.

Important requirements:
- `biz` is the HTTP entry layer.
- `internal/model` holds shared service-internal objects.
- `internal/service` holds business logic.
- `internal/dal/db` holds MySQL repositories.
- HTTP routing stays handwritten.

Concrete work:
- Move HTTP handlers/routing/validation into `services/media-service/biz/upload/*` and `services/media-service/biz/router.go`.
- Add `services/media-service/biz/model/upload.go` for HTTP DTOs if needed.
- Add `services/media-service/internal/model/file.go` for service-internal objects.
- Move application logic into `services/media-service/internal/service/*`.
- Move repository/storage init into `services/media-service/internal/dal/*`.
- Update `services/media-service/cmd/main.go` and `services/media-service/cmd/router.go` to wire the new structure.
- Keep existing tests green.

Validate with:
- `go test ./services/media-service/...`
- `go test ./...`

Report back:
- files changed
- tests run and result
- any concerns