# Microservices Runtime And Monolith Removal Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 让当前微服务骨架真正本地可运行、可通过 gateway 访问，并在验收通过后删除单体遗留代码。

**Architecture:** 先把共享 session、数据库迁移、本地启动、gateway 转发和服务间调用跑通，形成端到端闭环。再接入真实 Nacos 与 OpenTelemetry SDK，最后按 import 依赖证明逐批删除 `biz/`、`service/`、旧 IDL 和旧单体入口。删除单体前必须保留可重复的 smoke/e2e 验证脚本。

**Tech Stack:** Go、Hertz、Kitex-ready IDL、GORM、MySQL、Redis、Nacos、OpenTelemetry、Docker Compose、Kubernetes、PowerShell/Makefile

## Global Constraints

- 新增或改写的计划、运行手册、迁移说明统一使用中文。
- 外部 HTTP 流量最终只能进入 `gateway`，不能再直接访问单体路由。
- 登录鉴权使用 `session cookie + Redis`，不能切回 JWT。
- 每个服务拥有自己的数据库 schema 和迁移脚本。
- 一个服务不能直接读取另一个服务的数据库表。
- `article`、`category`、`tag` 保留在 `content-service` 一个内容域内。
- 公开文章详情按文章 `id` 查询。
- 删除单体遗留代码前必须通过 `go test ./...` 和端到端 smoke 验证。
- 删除时不能删除仍被新服务 import 的公共包。

---

## File Structure

- Modify: `pkg/xauth/session.go` - 从内存 session 切换为可插拔 store，新增 Redis store。
- Create: `pkg/xauth/redis_store.go` - Redis session 存取实现。
- Create: `pkg/xauth/redis_store_test.go` - 使用 fake Redis store 验证 session store 行为。
- Modify: `services/auth-service/internal/application/auth_service.go` - auth-service 使用共享 Redis session store。
- Modify: `services/auth-service/cmd/main.go` - 启动时初始化 Redis session store。
- Create: `pkg/xruntime/migrate.go` - 本地迁移执行器，按服务执行 SQL 文件。
- Create: `deploy/docker/init/001_databases.sql` - 创建 `auth_db`、`media_db`、`content_db`。
- Modify: `deploy/docker/compose.yaml` - 加入业务服务、Redis/MySQL 初始化和健康检查。
- Modify: `Makefile` - 新增 `micro-up`、`micro-down`、`micro-test`、`micro-smoke`。
- Modify: `services/gateway/internal/router/router.go` - 增加真实反向代理路由。
- Create: `services/gateway/internal/proxy/reverse_proxy.go` - Hertz 到下游 HTTP 服务的转发器。
- Create: `services/gateway/internal/proxy/reverse_proxy_test.go` - 验证路径重写与请求转发。
- Modify: `services/web-bff/internal/assembler/article_page.go` - 使用可配置 content-service 地址或发现结果。
- Create: `pkg/xnacos/nacos_client.go` - 接入真实 Nacos SDK 的实现点。
- Create: `pkg/xotel/otel_setup.go` - 接入真实 OpenTelemetry exporter 的实现点。
- Create: `scripts/smoke/microservices_smoke.ps1` - 本地端到端验收脚本。
- Delete after validation: `biz/`、`service/`、旧根目录单体入口和旧 IDL 中已经被 `idl/auth`、`idl/content`、`idl/media` 替代的文件。

### Task 1: Redis Session Store

**Files:**
- Modify: `pkg/xauth/session.go`
- Create: `pkg/xauth/redis_store.go`
- Create: `pkg/xauth/redis_store_test.go`
- Modify: `services/auth-service/internal/application/auth_service.go`
- Modify: `services/auth-service/cmd/main.go`

**Interfaces:**
- Consumes: `pkg/xauth.Claims`
- Produces: `type Store interface { Save(ctx context.Context, sessionID string, claims *Claims, ttl time.Duration) error; Get(ctx context.Context, sessionID string) (*Claims, error); Delete(ctx context.Context, sessionID string) error }`
- Produces: `func UseStore(store Store)`
- Produces: `func CreateSession(ctx context.Context, userID int64, username string, roles []string) (string, error)`
- Produces: `func ParseSession(ctx context.Context, sessionID string) (*Claims, error)`

- [ ] **Step 1: Write the failing store test**

```go
package xauth

import (
	"context"
	"testing"
	"time"
)

func TestSessionUsesConfiguredStore(t *testing.T) {
	store := newMemoryStoreForTest()
	UseStore(store)
	t.Cleanup(func() { UseStore(newMemoryStoreForTest()) })

	sessionID, err := CreateSession(context.Background(), 7, "admin", []string{"super_admin"})
	if err != nil {
		t.Fatalf("CreateSession returned error: %v", err)
	}

	claims, err := ParseSession(context.Background(), sessionID)
	if err != nil {
		t.Fatalf("ParseSession returned error: %v", err)
	}
	if claims.UserID != 7 {
		t.Fatalf("expected user id 7, got %d", claims.UserID)
	}
	if claims.ExpiresAt.Before(time.Now()) {
		t.Fatal("expected future expiration")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/xauth -run TestSessionUsesConfiguredStore -v`  
Expected: FAIL with `UseStore` or `CreateSession` signature mismatch.

- [ ] **Step 3: Implement minimal store abstraction**

```go
type Store interface {
	Save(ctx context.Context, sessionID string, claims *Claims, ttl time.Duration) error
	Get(ctx context.Context, sessionID string) (*Claims, error)
	Delete(ctx context.Context, sessionID string) error
}

var activeStore Store = newMemoryStoreForTest()

func UseStore(store Store) {
	if store == nil {
		activeStore = newMemoryStoreForTest()
		return
	}
	activeStore = store
}
```

- [ ] **Step 4: Add Redis-backed implementation**

`pkg/xauth/redis_store.go` must expose:

```go
func NewRedisStore(pool RedisPool, prefix string) Store
```

Use JSON encoding for `Claims`, store key format `session:<sessionID>`, and TTL from `SessionStoreConfig.ExpireHour`.

- [ ] **Step 5: Update auth-service call sites**

Change auth-service application methods to use context-aware calls:

```go
sessionID, err := xauth.CreateSession(ctx, userID, resolvedUsername, roles)
claims, err := xauth.ParseSession(ctx, sessionID)
```

- [ ] **Step 6: Verify**

Run: `go test ./pkg/xauth ./services/auth-service/... -v`  
Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add pkg/xauth services/auth-service
git commit -m "feat: use redis backed session store"
```

### Task 2: Local Runtime And Migrations

**Files:**
- Create: `pkg/xruntime/migrate.go`
- Create: `pkg/xruntime/migrate_test.go`
- Create: `deploy/docker/init/001_databases.sql`
- Modify: `deploy/docker/compose.yaml`
- Modify: `Makefile`
- Modify: `docs/runbooks/local-microservices.md`

**Interfaces:**
- Produces: `func MigrationFiles(root string) ([]string, error)`
- Produces: `make micro-up`
- Produces: `make micro-down`
- Produces: `make micro-test`

- [ ] **Step 1: Write failing migration file ordering test**

```go
package xruntime

import "testing"

func TestMigrationFilesReturnsSortedSQLFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "002_second.sql", "select 2;")
	writeFile(t, dir, "001_first.sql", "select 1;")
	writeFile(t, dir, "notes.md", "skip")

	files, err := MigrationFiles(dir)
	if err != nil {
		t.Fatalf("MigrationFiles returned error: %v", err)
	}
	if len(files) != 2 || files[0] == files[1] || !endsWith(files[0], "001_first.sql") {
		t.Fatalf("unexpected files: %#v", files)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/xruntime -run TestMigrationFilesReturnsSortedSQLFiles -v`  
Expected: FAIL because `MigrationFiles` is undefined.

- [ ] **Step 3: Implement migration file discovery**

```go
func MigrationFiles(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, filepath.Join(root, entry.Name()))
	}
	sort.Strings(files)
	return files, nil
}
```

- [ ] **Step 4: Make Compose actually boot dependencies and business services**

`deploy/docker/compose.yaml` must include:

- `mysql` with `deploy/docker/init/001_databases.sql`
- `redis`
- `nacos`
- `otel-collector`
- `auth-service`
- `media-service`
- `content-service`
- `web-bff`
- `gateway`

- [ ] **Step 5: Add Makefile commands**

```makefile
micro-up:
	docker compose -f deploy/docker/compose.yaml up -d --build

micro-down:
	docker compose -f deploy/docker/compose.yaml down

micro-test:
	go test ./...

micro-smoke:
	powershell -ExecutionPolicy Bypass -File scripts/smoke/microservices_smoke.ps1
```

- [ ] **Step 6: Verify**

Run: `go test ./pkg/xruntime -v`  
Expected: PASS.

Run: `docker compose -f deploy/docker/compose.yaml config`  
Expected: command exits 0.

- [ ] **Step 7: Commit**

```bash
git add pkg/xruntime deploy/docker Makefile docs/runbooks/local-microservices.md
git commit -m "feat: add local microservice runtime"
```

### Task 3: Gateway Reverse Proxy

**Files:**
- Create: `services/gateway/internal/proxy/reverse_proxy.go`
- Create: `services/gateway/internal/proxy/reverse_proxy_test.go`
- Modify: `services/gateway/internal/router/router.go`
- Modify: `services/gateway/cmd/main.go`

**Interfaces:**
- Produces: `func NewReverseProxy(targetBaseURL string, stripPrefix string) app.HandlerFunc`
- Produces routes:
  - `/api/auth/*path` -> `auth-service`
  - `/api/media/*path` -> `media-service`
  - `/api/content/*path` -> `content-service`
  - `/api/blog/*path` -> `web-bff`

- [ ] **Step 1: Write failing proxy path test**

```go
package proxy

import "testing"

func TestRewritePathStripsGatewayPrefix(t *testing.T) {
	got := RewritePath("/api/content/articles/12", "/api/content")
	if got != "/articles/12" {
		t.Fatalf("expected /articles/12, got %s", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/gateway/internal/proxy -run TestRewritePathStripsGatewayPrefix -v`  
Expected: FAIL because `RewritePath` is undefined.

- [ ] **Step 3: Implement path rewrite and proxy handler**

```go
func RewritePath(path string, stripPrefix string) string {
	rewritten := strings.TrimPrefix(path, stripPrefix)
	if rewritten == "" {
		return "/"
	}
	if !strings.HasPrefix(rewritten, "/") {
		return "/" + rewritten
	}
	return rewritten
}
```

- [ ] **Step 4: Register gateway routes**

`RegisterRoutes` must call:

```go
h.Any("/api/auth/*path", proxy.NewReverseProxy(deps.AuthBaseURL, "/api/auth"))
h.Any("/api/media/*path", proxy.NewReverseProxy(deps.MediaBaseURL, "/api/media"))
h.Any("/api/content/*path", proxy.NewReverseProxy(deps.ContentBaseURL, "/api/content"))
h.Any("/api/blog/*path", proxy.NewReverseProxy(deps.BFFBaseURL, "/api/blog"))
```

- [ ] **Step 5: Verify**

Run: `go test ./services/gateway/... -v`  
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add services/gateway
git commit -m "feat: route gateway traffic to services"
```

### Task 4: Real Service Discovery And Tracing Wiring

**Files:**
- Modify: `pkg/xnacos/client.go`
- Modify: `pkg/xotel/setup.go`
- Modify: `services/auth-service/cmd/main.go`
- Modify: `services/media-service/cmd/main.go`
- Modify: `services/content-service/cmd/main.go`
- Modify: `services/web-bff/cmd/main.go`
- Modify: `services/gateway/cmd/main.go`
- Modify: `docs/runbooks/k8s-deploy.md`

**Interfaces:**
- Produces: `func RegisterService(ctx context.Context, serviceName string, host string, port uint64) error`
- Produces: `func ResolveService(ctx context.Context, serviceName string) (string, error)`
- Produces: `func SetupTracerProvider(ctx context.Context, serviceName string, endpoint string) (func(context.Context) error, error)`

- [ ] **Step 1: Write failing Nacos validation test**

```go
package xnacos

import (
	"context"
	"testing"
)

func TestRegisterServiceRequiresServiceName(t *testing.T) {
	client := NewMemoryClient()
	err := client.RegisterService(context.Background(), "", "127.0.0.1", 9001)
	if err == nil {
		t.Fatal("expected error")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/xnacos -run TestRegisterServiceRequiresServiceName -v`  
Expected: FAIL because `NewMemoryClient` or `RegisterService` is undefined.

- [ ] **Step 3: Implement discovery interface before SDK binding**

```go
type Discovery interface {
	RegisterService(ctx context.Context, serviceName string, host string, port uint64) error
	ResolveService(ctx context.Context, serviceName string) (string, error)
}
```

- [ ] **Step 4: Bind command startup**

Every service `cmd/main.go` must:

```go
shutdown, err := xotel.SetupTracerProvider(context.Background(), serviceName, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
if err != nil {
	log.Fatal(err)
}
defer shutdown(context.Background())
```

- [ ] **Step 5: Verify**

Run: `go test ./pkg/xnacos ./pkg/xotel ./services/... -v`  
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add pkg/xnacos pkg/xotel services docs/runbooks/k8s-deploy.md
git commit -m "feat: wire discovery and tracing startup"
```

### Task 5: End-To-End Smoke Verification

**Files:**
- Create: `scripts/smoke/microservices_smoke.ps1`
- Create: `scripts/smoke/README.md`
- Modify: `Makefile`
- Modify: `docs/runbooks/local-microservices.md`

**Interfaces:**
- Produces: `scripts/smoke/microservices_smoke.ps1`
- Produces: `make micro-smoke`

- [ ] **Step 1: Write smoke script skeleton**

```powershell
$ErrorActionPreference = "Stop"

function Assert-StatusOk($Response, $Name) {
  if ($Response.StatusCode -lt 200 -or $Response.StatusCode -ge 300) {
    throw "$Name failed with status $($Response.StatusCode)"
  }
}
```

- [ ] **Step 2: Add health checks**

The script must call:

```powershell
Invoke-WebRequest "http://127.0.0.1:8888/healthz"
Invoke-WebRequest "http://127.0.0.1:9001/me" -SkipHttpErrorCheck
Invoke-WebRequest "http://127.0.0.1:9003/articles?page=1&page_size=1"
```

- [ ] **Step 3: Add login cookie flow**

The script must:

```powershell
$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$body = @{ username = "admin"; password = "admin" } | ConvertTo-Json
Invoke-WebRequest "http://127.0.0.1:9001/login" -Method Post -Body $body -ContentType "application/json" -WebSession $session
Invoke-WebRequest "http://127.0.0.1:9001/me" -WebSession $session
```

- [ ] **Step 4: Verify**

Run: `go test ./...`  
Expected: PASS.

Run after `make micro-up`: `make micro-smoke`  
Expected: script exits 0.

- [ ] **Step 5: Commit**

```bash
git add scripts/smoke Makefile docs/runbooks/local-microservices.md
git commit -m "test: add microservice smoke verification"
```

### Task 6: Remove Monolith Imports

**Files:**
- Modify or delete only after imports are gone:
  - `biz/`
  - `service/`
  - `dal/db/`
  - root monolith route wiring
  - legacy IDL files at `idl/article.thrift`, `idl/category.thrift`, `idl/tag.thrift`, `idl/upload.thrift`, `idl/rbac.thrift`, `idl/auth.thrift`

**Interfaces:**
- Consumes: all microservice packages from Tasks 1-5
- Produces: no imports from `github.com/Loe1210/personal-site/biz`
- Produces: no imports from `github.com/Loe1210/personal-site/service`
- Produces: no imports from `github.com/Loe1210/personal-site/dal/db` in service packages

- [ ] **Step 1: Prove monolith imports still exist**

Run: `rg -n "github.com/Loe1210/personal-site/(biz|service|dal/db)" --glob "*.go"`  
Expected: output lists remaining legacy imports.

- [ ] **Step 2: Move required shared types into service-owned packages**

If a new service imports `dal/db`, move the required model into that service's `internal/repository/mysql`. If a new service imports `biz/model`, define an application DTO in that service.

- [ ] **Step 3: Delete one legacy area at a time**

Delete in this order:

```bash
git rm -r biz
go test ./...
git commit -m "refactor: remove legacy biz handlers"

git rm -r service
go test ./...
git commit -m "refactor: remove legacy service layer"

git rm idl/article.thrift idl/category.thrift idl/tag.thrift idl/upload.thrift idl/rbac.thrift idl/auth.thrift
go test ./...
git commit -m "refactor: remove legacy thrift contracts"
```

- [ ] **Step 4: Remove or shrink `dal/db`**

Only remove `dal/db` after:

```bash
rg -n "dal/db" --glob "*.go"
```

Expected: no output from new service packages.

- [ ] **Step 5: Final verification**

Run: `go test ./...`  
Expected: PASS.

Run after `make micro-up`: `make micro-smoke`  
Expected: script exits 0.

- [ ] **Step 6: Commit**

```bash
git add .
git commit -m "refactor: remove monolith legacy code"
```

## Self-Review

- Spec coverage: covers Redis session, local runtime, gateway routing, service discovery/tracing, smoke verification, and monolith deletion.
- Placeholder scan: no open-ended implementation placeholders are intentionally left; every task includes exact files, functions, commands, and expected verification.
- Type consistency: session store, migration discovery, proxy routing, discovery, tracing, and smoke script interfaces are defined before later tasks consume them.

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-07-13-microservices-runtime-and-monolith-removal.md`. Two execution options:

**1. Subagent-Driven (recommended)** - dispatch a fresh subagent per task, review between tasks, fast iteration.

**2. Inline Execution** - execute tasks in this session using executing-plans, batch execution with checkpoints.

Which approach?
