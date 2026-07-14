# Remove Monolith Legacy Code Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在微服务本地运行和 smoke 验证已经通过的前提下，逐块删除旧单体入口、`biz/`、`service/`、`dal/db` 和旧根 IDL，保留可运行的微服务闭环。

**Architecture:** 先删除不再被微服务引用的旧入口和旧业务层，再删除旧数据访问层和旧 Thrift 文件。每一步都用 import 搜索、`go test ./...`、`make micro-smoke` 证明没有破坏现有微服务能力。删除过程只移除 legacy code，不把旧模型重新引回新服务。

**Tech Stack:** Go、Hertz、GORM、MySQL、Redis、Docker Compose、PowerShell、Makefile、ripgrep、Git

## Global Constraints

- 不删除 `services/`、`pkg/xauth`、`pkg/xotel`、`pkg/xnacos`、`pkg/xruntime`、`configs/`、`deploy/docker/compose.yaml`、`Dockerfile`。
- 删除旧单体前必须已经通过 `go test ./...` 和 `make micro-smoke`；当前基线已通过。
- 每个删除任务之后必须重新运行 `go test ./...`。
- 删除 `dal/db` 和旧 IDL 前必须确认新微服务没有 import `github.com/Loe1210/personal-site/dal/db`。
- 只提交本计划相关文件，不提交当前工作区里已有的无关文档改动：`docs/pending-tasks.md`、`docs/superpowers/specs/2026-07-13-phase-a-b-swagger-cleanup-design.md`、`.superpowers/`、旧未跟踪计划/规格文件。
- 如果任一验证失败，停止继续删除，先定位失败根因。

---

## File Structure

- Delete: `main.go` - 旧单体 Hertz 入口，已由 `services/gateway/cmd` 和各微服务入口替代。
- Delete: `biz/` - 旧单体 HTTP handler、route、request model。
- Delete: `service/` - 旧单体业务服务层。
- Delete: `pkg/middleware/session/rbac.go` - 旧单体 RBAC 中间件，依赖 `dal/db`，微服务网关当前不使用。
- Delete: `dal/db/` - 旧单体数据库模型、初始化、种子逻辑，已被各服务自己的 repository/migrate 替代。
- Delete: `idl/article.thrift`、`idl/category.thrift`、`idl/tag.thrift`、`idl/upload.thrift`、`idl/rbac.thrift`、`idl/auth.thrift` - 旧根 IDL，保留 `idl/auth/auth.thrift`、`idl/content/content.thrift`、`idl/media/media.thrift`。
- Modify: `README.md` only if a legacy monolith command still references root `main.go` after deletion.
- Modify: `docs/runbooks/local-microservices.md` only if smoke/runbook needs a note after deletion.

## Task 1: Delete Legacy Monolith Entrypoint And Biz Layer

**Files:**
- Delete: `main.go`
- Delete: `biz/`

**Interfaces:**
- Consumes: existing microservice commands under `services/*/cmd`
- Produces: no root single-binary application package

- [ ] **Step 1: Prove old entrypoint and biz imports exist**

Run:

```powershell
rg -n 'github.com/Loe1210/personal-site/(biz|service|dal/db)' --glob '*.go' main.go biz service pkg services dal
```

Expected: output includes `main.go` and `biz/` legacy imports.

- [ ] **Step 2: Delete root monolith entrypoint and biz layer**

Run:

```powershell
git rm -r main.go biz
```

Expected: Git stages deletion of `main.go` and every file under `biz/`.

- [ ] **Step 3: Verify remaining imports**

Run:

```powershell
rg -n 'github.com/Loe1210/personal-site/(biz|service|dal/db)' --glob '*.go' service pkg services dal
```

Expected: remaining output only comes from `service/` and `pkg/middleware/session/rbac.go`, not from `services/`.

- [ ] **Step 4: Run tests**

Run:

```powershell
go test ./...
```

Expected: PASS.

- [ ] **Step 5: Commit**

Run:

```powershell
git commit -m "refactor: remove legacy monolith handlers"
```

## Task 2: Delete Legacy Service Layer And RBAC Middleware

**Files:**
- Delete: `service/`
- Delete: `pkg/middleware/session/rbac.go`

**Interfaces:**
- Consumes: no new service imports from `service/`
- Produces: no import path `github.com/Loe1210/personal-site/service`

- [ ] **Step 1: Prove service layer is isolated**

Run:

```powershell
rg -n 'github.com/Loe1210/personal-site/service|github.com/Loe1210/personal-site/biz' --glob '*.go'
```

Expected: output only references `service/` files or no output after Task 1.

- [ ] **Step 2: Delete legacy service layer and old RBAC middleware**

Run:

```powershell
git rm -r service pkg/middleware/session/rbac.go
```

Expected: Git stages deletion of `service/` and `pkg/middleware/session/rbac.go`.

- [ ] **Step 3: Verify no biz/service imports remain**

Run:

```powershell
rg -n 'github.com/Loe1210/personal-site/(biz|service)' --glob '*.go'
```

Expected: no output.

- [ ] **Step 4: Run tests**

Run:

```powershell
go test ./...
```

Expected: PASS.

- [ ] **Step 5: Commit**

Run:

```powershell
git commit -m "refactor: remove legacy service layer"
```

## Task 3: Delete Legacy DAL And Root IDL Contracts

**Files:**
- Delete: `dal/db/`
- Delete: `idl/article.thrift`
- Delete: `idl/category.thrift`
- Delete: `idl/tag.thrift`
- Delete: `idl/upload.thrift`
- Delete: `idl/rbac.thrift`
- Delete: `idl/auth.thrift`

**Interfaces:**
- Consumes: service-owned migrations under `services/*/internal/repository/mysql/migrate.go`
- Produces: no import path `github.com/Loe1210/personal-site/dal/db`
- Preserves: `idl/auth/auth.thrift`, `idl/content/content.thrift`, `idl/media/media.thrift`

- [ ] **Step 1: Prove new services do not import old DAL**

Run:

```powershell
rg -n 'github.com/Loe1210/personal-site/dal/db|dal/db' --glob '*.go' services pkg configs
```

Expected: no output from `services/`; if output exists only from already-deleted legacy files, stop and inspect.

- [ ] **Step 2: Delete old DAL and root IDL files**

Run:

```powershell
git rm -r dal/db
git rm idl/article.thrift idl/category.thrift idl/tag.thrift idl/upload.thrift idl/rbac.thrift idl/auth.thrift
```

Expected: Git stages deletion of old DAL and root IDL files while preserving nested service IDL directories.

- [ ] **Step 3: Verify old DAL imports are gone**

Run:

```powershell
rg -n 'github.com/Loe1210/personal-site/dal/db|dal/db' --glob '*.go'
```

Expected: no output.

- [ ] **Step 4: Run tests**

Run:

```powershell
go test ./...
```

Expected: PASS.

- [ ] **Step 5: Commit**

Run:

```powershell
git commit -m "refactor: remove legacy dal and thrift contracts"
```

## Task 4: Final Runtime Verification And Documentation Cleanup

**Files:**
- Modify: `README.md` if it still describes `go run .` as the primary entrypoint.
- Modify: `docs/runbooks/local-microservices.md` only if needed to clarify microservice-only startup.

**Interfaces:**
- Produces: verified microservice-only runtime.

- [ ] **Step 1: Scan docs for legacy startup instructions**

Run:

```powershell
rg -n 'go run \.|docker compose up --build|docker-compose.yml|单体|monolith' README.md docs
```

Expected: identify only historical notes or sections that need wording updates.

- [ ] **Step 2: Update README if needed**

If README still says root `go run .` is the main local run mode, replace that paragraph with:

```markdown
## 本地运行

当前推荐使用微服务版 Docker Compose：

```bash
make micro-up
make micro-smoke
```

旧单体入口已删除，外部 HTTP 流量统一从 gateway 进入：`http://localhost:8888`。
```

- [ ] **Step 3: Run full test and smoke verification**

Run:

```powershell
go test ./...
C:\ProgramData\chocolatey\bin\make.exe micro-smoke
```

Expected: both PASS.

- [ ] **Step 4: Verify old paths are absent**

Run:

```powershell
Test-Path main.go
Test-Path biz
Test-Path service
Test-Path dal\db
rg -n 'github.com/Loe1210/personal-site/(biz|service|dal/db)' --glob '*.go'
```

Expected: all `Test-Path` results are `False`; final `rg` has no output.

- [ ] **Step 5: Commit documentation cleanup if files changed**

Run:

```powershell
git add README.md docs/runbooks/local-microservices.md
git commit -m "docs: document microservice-only runtime"
```

Only run this commit if `git diff -- README.md docs/runbooks/local-microservices.md` has changes.

## Self-Review

- Spec coverage: covers proof of legacy imports, deletion of old entrypoint, `biz/`, `service/`, RBAC middleware, `dal/db`, root IDL, and final smoke verification.
- Placeholder scan: no TBD/TODO placeholders; every task has exact commands and expected result.
- Type consistency: no new Go interfaces are introduced; deletion tasks depend only on verified absence of imports and existing microservice runtime.

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-07-14-remove-monolith-legacy-code.md`. Execution will use Inline Execution with `superpowers:executing-plans` because the user requested planning and execution in this session.