# 微服务服务内部分层统一实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把现有 `auth-service`、`content-service`、`media-service` 统一整理成 `biz / cmd / idl / internal / pkg` 的微服务目录风格，并保留当前可运行的 HTTP + RPC 闭环。

**Architecture:** `biz` 只放 HTTP 入口层和路由骨架，`internal/model` 只放服务内部共享对象，`internal/service` 放业务逻辑，`internal/dal` 放数据库、缓存、RPC 客户端和外部访问适配。HTTP 路由继续手写，RPC 契约保持 IDL 先行、生成物只读。

**Tech Stack:** Go, Hertz, Kitex, Thrift, MySQL, Redis, Makefile

## Global Constraints

- HTTP 路由继续手写，`biz` 只承载 HTTP 入口层，不把核心业务逻辑塞回 handler。
- `internal/model` 是每个服务的内部共享对象边界，不再让 `handler`、`service`、`repository` 直接互相暴露结构体。
- RPC 接口统一先定义在 `idl/`，生成代码只读，禁止手工修改生成文件。
- 字段序号一旦定义不可随意变更。
- 登录鉴权继续使用 session + cookie + Redis，不切回 JWT。
- 已删除的单体入口、`biz/`、`service/`、`dal/db/` 和旧根 IDL 不允许恢复。
- `go test ./...` 和 `make micro-smoke` 必须继续通过。

---

### Task 1: 统一 content-service 的分层结构

**Files:**
- Create: `services/content-service/biz/router.go`
- Create: `services/content-service/biz/article/handler.go`
- Create: `services/content-service/biz/article/route.go`
- Create: `services/content-service/biz/article/validate.go`
- Create: `services/content-service/biz/model/article.go`
- Create: `services/content-service/internal/model/article.go`
- Create: `services/content-service/internal/service/article_service.go`
- Create: `services/content-service/internal/service/tag_service.go`
- Create: `services/content-service/internal/service/category_service.go`
- Create: `services/content-service/internal/dal/init.go`
- Create: `services/content-service/internal/dal/db/article_repository.go`
- Create: `services/content-service/internal/dal/db/tag_repository.go`
- Create: `services/content-service/internal/dal/db/category_repository.go`
- Modify: `services/content-service/cmd/main.go`
- Modify: `services/content-service/cmd/router.go`
- Modify: `services/content-service/internal/handler/http/article.go`
- Modify: `services/content-service/internal/handler/rpc/content.go`
- Modify: `services/content-service/internal/application/article_service.go`
- Modify: `services/content-service/internal/application/tag_service.go`
- Modify: `services/content-service/internal/application/category_service.go`
- Modify: `services/content-service/internal/repository/mysql/article_repository.go`

**Interfaces:**
- Consumes: `internal/model.ArticleDetail`, `internal/model.TagDTO`, `internal/model.ListFilter`, `internal/model.ListResult`.
- Produces: `biz/article.Handler` for HTTP, `internal/service.ArticleService` for use cases, `internal/dal/db.ArticleRepository` for persistence.

- [ ] **Step 1: 先把现有 content-service 的调用链梳理成新边界**

确认 article 的 HTTP 入参、service 入参、repository 返回值分别落在哪一层，要求最终只有 `internal/model` 被 service 与 repository 共享，`biz` 只接收和返回 HTTP 结构。

- [ ] **Step 2: 迁移 HTTP 层到 `biz/article`**

把当前 `internal/handler/http/article.go` 的路由、参数校验、响应封装拆成 `biz/article/handler.go`、`route.go`、`validate.go`，并让 `cmd/router.go` 只负责挂载业务路由。

- [ ] **Step 3: 迁移业务逻辑到 `internal/service`**

把 `internal/application/article_service.go`、`tag_service.go`、`category_service.go` 改成 `internal/service/*`，方法签名只依赖 `internal/model`，不再暴露 HTTP 层类型。

- [ ] **Step 4: 迁移数据访问到 `internal/dal/db`**

把 `internal/repository/mysql/article_repository.go` 迁到 `internal/dal/db/article_repository.go`，把 MySQL 初始化收口到 `internal/dal/init.go`，保持 article/tag/category 的读写行为不变。

- [ ] **Step 5: 跑 content-service 的测试和最小烟雾验证**

Run: `go test ./services/content-service/...`
Expected: 通过，且没有残留旧路径 import。

Run: `make micro-smoke`
Expected: gateway health、auth `/me`、content article list、login cookie flow 全部通过。

### Task 2: 统一 auth-service 的分层结构

**Files:**
- Create: `services/auth-service/biz/router.go`
- Create: `services/auth-service/biz/authenticator/handler.go`
- Create: `services/auth-service/biz/authenticator/route.go`
- Create: `services/auth-service/biz/authenticator/validate.go`
- Create: `services/auth-service/biz/model/auth.go`
- Create: `services/auth-service/internal/model/user.go`
- Create: `services/auth-service/internal/service/auth_service.go`
- Create: `services/auth-service/internal/dal/init.go`
- Create: `services/auth-service/internal/dal/db/user_repository.go`
- Modify: `services/auth-service/cmd/main.go`
- Modify: `services/auth-service/cmd/router.go`
- Modify: `services/auth-service/internal/handler/http/login.go`
- Modify: `services/auth-service/internal/handler/rpc/auth.go`
- Modify: `services/auth-service/internal/application/auth_service.go`
- Modify: `services/auth-service/internal/repository/mysql/user_repository.go`

**Interfaces:**
- Consumes: `internal/model.User`, `internal/model.SessionBundle`, `internal/model.AuthContext`。
- Produces: `biz/authenticator.Handler` for HTTP login/logout/me, `internal/service.AuthService` for session and permission logic, `internal/dal/db.UserRepository` for user lookup.

- [ ] **Step 1: 把 session 登录对象收进 `internal/model`**

确保 `User`、`SessionBundle`、`AuthContext` 只在服务内部流转，HTTP 层只拿到 `biz/model` 的登录请求和响应结构。

- [ ] **Step 2: 迁移 HTTP 登录路由到 `biz/authenticator`**

把 `/login`、`/logout`、`/me` 的路由注册和请求校验拆到 `biz/authenticator`，保留 session cookie 写入逻辑和 Redis 会话存储。

- [ ] **Step 3: 迁移 auth 业务逻辑到 `internal/service`**

把鉴权、会话校验、当前用户查询和权限校验集中到 `internal/service/auth_service.go`，接口不再返回 HTTP 结构体。

- [ ] **Step 4: 迁移用户仓储到 `internal/dal/db`**

把用户查询、角色加载、权限判断实现放进 `internal/dal/db/user_repository.go`，保留当前数据库表和密码校验行为。

- [ ] **Step 5: 跑 auth-service 测试**

Run: `go test ./services/auth-service/...`
Expected: 通过，且 `session cookie redis` 的登录链路仍然可编译。

### Task 3: 统一 media-service 的分层结构

**Files:**
- Create: `services/media-service/biz/router.go`
- Create: `services/media-service/biz/upload/handler.go`
- Create: `services/media-service/biz/upload/route.go`
- Create: `services/media-service/biz/upload/validate.go`
- Create: `services/media-service/biz/model/upload.go`
- Create: `services/media-service/internal/model/file.go`
- Create: `services/media-service/internal/service/media_service.go`
- Create: `services/media-service/internal/dal/init.go`
- Create: `services/media-service/internal/dal/storage/local.go`
- Create: `services/media-service/internal/dal/db/file_repository.go`
- Modify: `services/media-service/cmd/main.go`
- Modify: `services/media-service/cmd/router.go`
- Modify: `services/media-service/internal/handler/http/upload.go`
- Modify: `services/media-service/internal/application/media_service.go`
- Modify: `services/media-service/internal/repository/mysql/file_repository.go`
- Modify: `services/media-service/internal/infra/storage/local.go`

**Interfaces:**
- Consumes: `internal/model.FileRecord`, `internal/model.UploadInput`。
- Produces: `biz/upload.Handler` for HTTP upload/get-file, `internal/service.MediaService` for upload flow, `internal/dal/db.FileRepository` for file metadata persistence.

- [ ] **Step 1: 把文件元数据统一进 `internal/model`**

确保上传入参、文件记录和数据库模型之间有清晰映射，不再让 HTTP 层直接依赖 repository 内部结构。

- [ ] **Step 2: 迁移上传路由到 `biz/upload`**

把 `POST /upload` 和 `GET /files/:id` 的路由与校验拆成独立的 `biz/upload` 包，保留现有上传行为和返回格式。

- [ ] **Step 3: 迁移 media 业务逻辑到 `internal/service`**

把本地存储选择、BizType 归一化、文件落库逻辑收口到 `internal/service/media_service.go`。

- [ ] **Step 4: 迁移文件仓储和存储适配到 `internal/dal`**

把 MySQL 记录保存和本地文件存储初始化收口到 `internal/dal/db` 与 `internal/dal/storage`。

- [ ] **Step 5: 跑 media-service 测试**

Run: `go test ./services/media-service/...`
Expected: 通过，且上传接口仍可正常编译。

### Task 4: 收口 IDL、文档和统一验证

**Files:**
- Modify: `Makefile`
- Modify: `README.md`
- Modify: `docs/rpc-development-guidelines.md`
- Modify: `docs/pending-tasks.md`
- Create: `idl/content/content.thrift` 或现有 `idl` 对应文件的更新版本
- Create: `idl/auth/auth.thrift` 或现有 `idl` 对应文件的更新版本
- Create: `idl/media/media.thrift` 或现有 `idl` 对应文件的更新版本

**Interfaces:**
- Consumes: 各服务的 HTTP handler、RPC handler、IDL 定义。
- Produces: 明确的生成命令、只读生成物约束、统一的服务启动和验证文档。

- [ ] **Step 1: 把 RPC 开发规范写进仓库文档**

补充 IDL 先行、字段序号固定、生成物只读、HTTP 手写的规范，避免后续又回到手改生成文件。

- [ ] **Step 2: 把 Makefile 的服务启动和生成命令写清楚**

把微服务启动、测试和后续 IDL 生成命令统一在 `Makefile`，让开发者只看一处入口就能知道怎么跑。

- [ ] **Step 3: 统一跑全量测试和 smoke**

Run: `go test ./...`
Expected: 全仓库通过。

Run: `make micro-smoke`
Expected: 微服务闭环通过。

- [ ] **Step 4: 清理旧路径和过渡说明**

确认没有遗留的单体旧目录引用后，再更新 README 和待办文档，避免新人误读当前架构。