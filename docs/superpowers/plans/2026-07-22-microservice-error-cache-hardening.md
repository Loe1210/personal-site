# Microservice Error Cache Hardening Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Unify expected error responses, harden auth/content RPC error semantics, add content-service-only local -> Redis -> MySQL article cache, expand smoke coverage, and document the architecture.

**Architecture:** Add a small shared base proto and Go response/error packages. Keep gateway as a reverse proxy plus auth gate. Put all content caching below content-service handlers in `ArticleService`.

**Tech Stack:** Go, Hertz, Kitex protobuf, redigo Redis, GORM MySQL, PowerShell smoke tests.

## Global Constraints

- Gateway stays thin and must not add content business logic.
- Three-layer content cache lives only in content-service.
- Tests are written before production implementation.
- Each phase runs `go test ./...`, `make proto-check`, and `make micro-smoke`.
- Commit after completion; do not deploy.

---

### Task 1: Shared Errors, Envelopes, And Recover

**Files:**
- Create: `idl/base/base.proto`
- Create: `internal/xerrors/errors.go`
- Create: `internal/xhttp/response.go`
- Create: `internal/xhttp/recover.go`
- Create: `internal/xsafe/goroutine.go`
- Modify: `idl/auth/auth.proto`
- Modify: `idl/content/content.proto`
- Modify: generated `kitex_gen` protobuf files after proto generation

**Interfaces:**
- Produces: `xerrors.AppError`, `xerrors.CodeOf(error)`, `xhttp.OK`, `xhttp.Fail`, `xhttp.Recover()`, `xsafe.Go`, `xsafe.DeferRecover`.

- [ ] Write failing tests for app error mapping, HTTP envelope shape, recover envelope, and gopool/goroutine recover helper.
- [ ] Update proto to include `base.BaseResp` in auth/content responses.
- [ ] Generate or mechanically update protobuf code.
- [ ] Implement minimal shared packages.
- [ ] Run phase verification.

### Task 2: Auth RPC Error Mapping

**Files:**
- Modify: `services/auth-service/internal/service/auth_service.go`
- Modify: `services/auth-service/internal/handler/rpc/auth.go`
- Modify: `services/gateway/internal/client/auth/kitex_client.go`
- Modify: `services/gateway/internal/middleware/auth.go`

**Interfaces:**
- Consumes: `base.BaseResp`, `xerrors.AppError`, `xhttp.Fail`.
- Produces: auth business errors that remain envelope-based at HTTP boundaries.

- [ ] Write failing tests for invalid session returning `base_resp` without RPC error.
- [ ] Write failing tests for gateway auth client mapping `base_resp` to typed errors.
- [ ] Write failing tests for middleware returning HTTP 200 envelope for auth business failures.
- [ ] Implement minimal service, handler, client, and middleware changes.
- [ ] Run phase verification.

### Task 3: Content Handler Errors And Cache

**Files:**
- Create: `services/content-service/internal/service/article_cache.go`
- Modify: `services/content-service/internal/service/article_service.go`
- Modify: `services/content-service/biz/article/handler.go`
- Modify: `services/content-service/cmd/main.go`

**Interfaces:**
- Consumes: `xerrors.AppError`, `xhttp.OK`, `xhttp.Fail`.
- Produces: `NewArticleServiceWithCache(repo ArticleGetter, cache ArticleCache)`.

- [ ] Write failing tests for handler parameter errors returning HTTP 200 envelope.
- [ ] Write failing tests for service not found returning a business error.
- [ ] Write failing tests proving local cache hit avoids Redis/MySQL, Redis hit backfills local, and MySQL hit backfills both.
- [ ] Implement minimal cache interfaces and service flow.
- [ ] Invalidate cache on create, update, and delete.
- [ ] Run phase verification.

### Task 4: Smoke And Docs

**Files:**
- Modify: `scripts/smoke/microservices_smoke.ps1`
- Modify: `docs/devlog/phase-17.md`
- Modify: `docs/architecture/microservices-current-state.md`

- [ ] Write smoke assertions for `/api/content/articles` and `/api/articles` 404.
- [ ] Update docs with error ownership, response envelope, recover, and cache boundary.
- [ ] Run final verification.
- [ ] Stage and commit.
