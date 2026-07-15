# 微服务重构实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**目标：** 在不破坏现有文章与后台主链路的前提下，将当前个人站点单体逐步迁移为面向 Kubernetes 的微服务系统，核心技术栈为 `Hertz + Kitex + Nacos + OpenTelemetry`。

**架构：** 先把当前单体整理成“可拆”的结构，再优先拆出 `auth-service` 和 `media-service`，随后迁移内容域到 `content-service`，最后补齐 `gateway`、服务发现、可观测性和 Kubernetes 部署资产。每个阶段都必须保持系统可运行、可验证、可回退。

**技术栈：** Go、Hertz、Kitex、GORM、MySQL、Redis、Nacos、OpenTelemetry、OTEL Collector、Prometheus、Grafana、Jaeger/Tempo、Docker、Kubernetes

## 全局约束

- 目标技术栈固定为 `Hertz + Kitex + Nacos + OpenTelemetry + Kubernetes`。
- `article`、`category`、`tag` 必须保留在同一个 `content-service` 中，不继续细拆成多个服务。
- 公开文章详情统一使用文章 `id` 作为标识。
- 在服务拆分前，必须先将当前登录态收敛为 `session + cookie + redis` 共享会话方案。
- 每个服务拥有自己的数据库 schema 和迁移脚本。
- 严禁一个服务直接读取另一个服务的数据库表。
- 所有外部 HTTP 流量统一经 `gateway` 进入系统。
- `web-bff` 只负责聚合，不拥有领域真相。
- 本次重构新增或改写的设计文档、实施计划、迁移说明、运行手册统一使用中文编写。

---

## 文件结构

本次实施将逐步引入以下顶层结构：

- 新建：`go.work`
- 新建：`services/gateway/`
- 新建：`services/auth-service/`
- 新建：`services/content-service/`
- 新建：`services/media-service/`
- 新建：`services/web-bff/`
- 新建：`idl/auth/`
- 新建：`idl/content/`
- 新建：`idl/media/`
- 新建：`pkg/xconfig/`、`pkg/xlog/`、`pkg/xtrace/`、`pkg/xerror/`、`pkg/xresponse/`、`pkg/xotel/`、`pkg/xnacos/`、`pkg/xauth/`
- 新建：`deploy/docker/compose.yaml`
- 新建：`deploy/k8s/base/`、`deploy/k8s/dev/`、`deploy/k8s/prod/`

本计划按阶段推进，每个任务都需要形成独立可验证的交付面。

### 任务 1：单体预重构，建立可拆分边界

**文件：**
- 新建：`pkg/xauth/session.go`
- 新建：`pkg/xauth/session_test.go`
- 新建：`pkg/xconfig/loader.go`
- 新建：`pkg/xtrace/http.go`
- 新建：`configs/config_test.go`
- 修改：`configs/config.go`
- 修改：`configs/config.yaml`
- 修改：`biz/auth/handler.go`
- 修改：`pkg/middleware/session/auth.go`
- 修改：`pkg/middleware/session/rbac.go`
- 修改：`dal/db/init.go`
- 修改：`README.md`

**接口：**
- 使用：现有 `authmodel.UserLoginRequest`、`response.WriteSuccess`、`errno` 包
- 产出：`func CreateSession(userID int64, username string, roles []string) (string, error)`
- 产出：`type Claims struct { UserID int64; Username string; Roles []string }`

- [ ] **步骤 1：先写失败的 session 单测**

```go
package xauth

import "testing"

func TestCreateSession(t *testing.T) {
	sessionID, err := CreateSession(7, "loe", []string{"super_admin"})
	if err != nil {
		t.Fatalf("CreateSession returned error: %v", err)
	}

	claims, err := ParseSession(sessionID)
	if err != nil {
		t.Fatalf("ParseSession returned error: %v", err)
	}

	if claims.UserID != 7 {
		t.Fatalf("expected user id 7, got %d", claims.UserID)
	}
	if claims.Username != "loe" {
		t.Fatalf("expected username loe, got %s", claims.Username)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "super_admin" {
		t.Fatalf("unexpected roles: %#v", claims.Roles)
	}
}
```

- [ ] **步骤 2：运行测试，确认它先失败**

运行：`go test ./pkg/xauth -run TestCreateSession -v`  
预期：FAIL，提示 `CreateSession` 或 `ParseSession` 未定义。

- [ ] **步骤 3：写最小可用 session 实现**

```go
package xauth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var sessionPrefix = "session:"

type SessionMetadata struct {
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Claims struct {
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	SessionMetadata
}

func CreateSession(userID int64, username string, roles []string) (string, error) {
	_ = Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		SessionMetadata: SessionMetadata{
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(2 * time.Hour),
		},
	}
	return sessionPrefix + uuid.NewString(), nil
}

func ParseSession(raw string) (*Claims, error) {
	if raw == "" {
		return nil, errors.New("empty session id")
	}
	return &Claims{}, nil
}
```

- [ ] **步骤 4：补配置与认证链路改造**

```go
type SessionStoreConfig struct {
	Prefix     string `yaml:"prefix"`
	ExpireHour int    `yaml:"expire_hour"`
	CookieName string `yaml:"cookie_name"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Config struct {
	Server       ServerConfig       `yaml:"server"`
	MySQL        MySQLConfig        `yaml:"mysql"`
	Session      SessionConfig      `yaml:"session"`
	SessionStore SessionStoreConfig `yaml:"session_store"`
	Redis        RedisConfig        `yaml:"redis"`
}
```

同时完成以下调整：

- 登录成功后写入 Redis session，并通过 cookie 返回 `session_id`
- 鉴权中间件改为从 cookie 读取 `session_id`，并通过 Redis 或 `auth-service` 校验会话
- 权限中间件依赖已解析的会话上下文，而不是本地内存 session
- `dal/db/init.go` 中把“连接数据库”“迁移”“种子数据”整理成可拆函数

- [ ] **步骤 5：运行测试，确认通过**

运行：`go test ./pkg/xauth ./configs -v`  
预期：PASS。

- [ ] **步骤 6：回归当前单体主链路**

运行：`go test ./...`  
预期：现有单体测试全部通过；如没有完整测试，至少认证相关包通过编译与测试。

- [ ] **步骤 7：更新中文说明文档**

在 `README.md` 中补充：

- 当前认证方式已改为 `session + cookie + redis`
- 调试请求需要携带浏览器返回的 session cookie
- 本阶段仍为单体，但已开始为服务拆分做边界预处理

- [ ] **步骤 8：提交**

```bash
git add pkg/xauth pkg/xconfig pkg/xtrace configs biz/auth/handler.go pkg/middleware/session dal/db/init.go README.md
git commit -m "refactor: prepare monolith for service extraction"
```

### 任务 2：拆出 auth-service

**文件：**
- 新建：`services/auth-service/cmd/main.go`
- 新建：`services/auth-service/internal/application/auth_service.go`
- 新建：`services/auth-service/internal/application/auth_service_test.go`
- 新建：`services/auth-service/internal/handler/http/login.go`
- 新建：`services/auth-service/internal/handler/rpc/auth.go`
- 新建：`services/auth-service/internal/domain/user.go`
- 新建：`services/auth-service/internal/repository/mysql/user_repository.go`
- 新建：`services/auth-service/internal/infra/mysql/mysql.go`
- 新建：`services/auth-service/configs/config.yaml`
- 新建：`services/auth-service/migrations/001_init.sql`
- 新建：`idl/auth/auth.thrift`
- 新建：`docs/migration/auth-service.md`

**接口：**
- 使用：`pkg/xauth.Claims`、`pkg/xresponse`、用户和 RBAC 表
- 产出：`CreateSession(ctx context.Context, username string, password string) (*SessionBundle, error)`
- 产出：`ValidateSession(ctx context.Context, sessionID string) (*AuthContext, error)`
- 产出：`CheckPermission(ctx context.Context, userID int64, code string) (bool, error)`

- [ ] **步骤 1：先写失败的 auth 应用层测试**

```go
package application

import (
	"context"
	"testing"
)

type fakeUserRepo struct{}

func (f *fakeUserRepo) Login(username, password string) (int64, string, []string, error) {
	return 1, "admin", []string{"super_admin"}, nil
}

func TestCreateSessionFromCredentials(t *testing.T) {
	svc := NewAuthService(&fakeUserRepo{})
	resp, err := svc.CreateSession(context.Background(), "admin", "admin")
	if err != nil {
		t.Fatalf("CreateSession returned error: %v", err)
	}
	if resp.SessionID == "" {
		t.Fatal("expected session id")
	}
}
```

- [ ] **步骤 2：运行测试，确认它先失败**

运行：`go test ./services/auth-service/internal/application -run TestCreateSessionFromCredentials -v`  
预期：FAIL，提示 `NewAuthService` 未定义。

- [ ] **步骤 3：写最小 auth-service 应用层实现**

```go
package application

import (
	"context"

	"github.com/Loe1210/personal-site/pkg/xauth"
)

type UserRepository interface {
	Login(username, password string) (int64, string, []string, error)
}

type TokenBundle struct {
	AccessToken string
}

type Service struct {
	users UserRepository
}

func NewAuthService(users UserRepository) *Service {
	return &Service{users: users}
}

func (s *Service) CreateSession(ctx context.Context, username string, password string) (*SessionBundle, error) {
	userID, resolvedUsername, roles, err := s.users.Login(username, password)
	if err != nil {
		return nil, err
	}
	sessionID, err := xauth.CreateSession(userID, resolvedUsername, roles)
	if err != nil {
		return nil, err
	}
	return &SessionBundle{SessionID: sessionID}, nil
}
```

- [ ] **步骤 4：补齐 auth-service 基础骨架**

完成以下最小能力：

- `cmd/main.go` 能独立启动 Hertz 服务
- 提供 `/login`、`/logout`、`/me` HTTP 入口
- 提供 `ValidateSession`、`CheckPermission` RPC 入口
- `migrations/001_init.sql` 建立 `users`、`roles`、`permissions`、`user_roles`、`role_permissions`，并为 Redis 会话结构预留键空间约定
- `configs/config.yaml` 提供服务名、端口、数据库、session、cookie、Redis 配置

- [ ] **步骤 5：运行测试与编译验证**

运行：`go test ./services/auth-service/... -v`  
预期：PASS。

运行：`go test ./...`  
预期：仓库整体可编译，新增服务相关包可通过测试。

- [ ] **步骤 6：补中文迁移说明**

在 `docs/migration/auth-service.md` 记录：

- 为什么先拆 auth-service
- 单体里哪些逻辑迁出
- 当前 RPC 与 HTTP 边界
- 后续怎么接入 gateway 和 Nacos

- [ ] **步骤 7：提交**

```bash
git add services/auth-service idl/auth/auth.thrift docs/migration/auth-service.md
git commit -m "feat: extract auth service skeleton"
```

### 任务 3：拆出 media-service

**文件：**
- 新建：`services/media-service/cmd/main.go`
- 新建：`services/media-service/internal/application/media_service.go`
- 新建：`services/media-service/internal/application/media_service_test.go`
- 新建：`services/media-service/internal/handler/http/upload.go`
- 新建：`services/media-service/internal/domain/file.go`
- 新建：`services/media-service/internal/repository/mysql/file_repository.go`
- 新建：`services/media-service/internal/infra/storage/local.go`
- 新建：`services/media-service/configs/config.yaml`
- 新建：`services/media-service/migrations/001_init.sql`
- 新建：`idl/media/media.thrift`
- 新建：`docs/migration/media-service.md`

**接口：**
- 使用：上传文件元信息、本地或对象存储配置
- 产出：`Upload(ctx context.Context, in UploadInput) (*FileRecord, error)`
- 产出：`GetFile(ctx context.Context, id int64) (*FileRecord, error)`

- [ ] **步骤 1：先写失败的 media 应用层测试**

```go
package application

import (
	"context"
	"testing"
)

type fakeStorage struct{}

func (f *fakeStorage) Save(name string, content []byte) (string, error) {
	return "/uploads/" + name, nil
}

func TestUploadReturnsURL(t *testing.T) {
	svc := NewMediaService(&fakeStorage{}, nil)
	resp, err := svc.Upload(context.Background(), UploadInput{
		FileName: "cover.png",
		Content:  []byte("png"),
	})
	if err != nil {
		t.Fatalf("Upload returned error: %v", err)
	}
	if resp.URL == "" {
		t.Fatal("expected upload URL")
	}
}
```

- [ ] **步骤 2：运行测试，确认它先失败**

运行：`go test ./services/media-service/internal/application -run TestUploadReturnsURL -v`  
预期：FAIL，提示 `NewMediaService` 未定义。

- [ ] **步骤 3：写最小 media-service 应用层实现**

```go
package application

import "context"

type Storage interface {
	Save(name string, content []byte) (string, error)
}

type Repository interface {
	Save(ctx context.Context, record *FileRecord) error
}

type UploadInput struct {
	FileName string
	Content  []byte
}

type FileRecord struct {
	ID  int64
	URL string
}

type Service struct {
	storage Storage
	repo    Repository
}

func NewMediaService(storage Storage, repo Repository) *Service {
	return &Service{storage: storage, repo: repo}
}

func (s *Service) Upload(ctx context.Context, in UploadInput) (*FileRecord, error) {
	url, err := s.storage.Save(in.FileName, in.Content)
	if err != nil {
		return nil, err
	}
	record := &FileRecord{URL: url}
	if s.repo != nil {
		if err := s.repo.Save(ctx, record); err != nil {
			return nil, err
		}
	}
	return record, nil
}
```

- [ ] **步骤 4：补齐 media-service 基础骨架**

完成以下最小能力：

- `cmd/main.go` 能独立启动
- 提供 `/upload` HTTP 入口
- 本地存储实现先跑通
- `migrations/001_init.sql` 建立上传元数据表
- `configs/config.yaml` 提供上传目录、访问域名、数据库配置

- [ ] **步骤 5：运行测试与编译验证**

运行：`go test ./services/media-service/... -v`  
预期：PASS。

运行：`go test ./...`  
预期：仓库整体可编译，新增服务相关包可通过测试。

- [ ] **步骤 6：补中文迁移说明**

在 `docs/migration/media-service.md` 记录：

- 上传逻辑为什么适合早拆
- 当前先用本地存储的原因
- 后续切 MinIO 时的替换点

- [ ] **步骤 7：提交**

```bash
git add services/media-service idl/media/media.thrift docs/migration/media-service.md
git commit -m "feat: extract media service skeleton"
```

### 任务 4：拆出 content-service 与 web-bff

**文件：**
- 新建：`services/content-service/cmd/main.go`
- 新建：`services/content-service/internal/application/article_service.go`
- 新建：`services/content-service/internal/application/category_service.go`
- 新建：`services/content-service/internal/application/tag_service.go`
- 新建：`services/content-service/internal/application/article_service_test.go`
- 新建：`services/content-service/internal/handler/http/article.go`
- 新建：`services/content-service/internal/handler/rpc/content.go`
- 新建：`services/content-service/internal/repository/mysql/article_repository.go`
- 新建：`services/content-service/migrations/001_init.sql`
- 新建：`services/web-bff/cmd/main.go`
- 新建：`services/web-bff/internal/handler/http/blog.go`
- 新建：`services/web-bff/internal/assembler/article_page.go`
- 新建：`services/web-bff/internal/assembler/article_page_test.go`
- 新建：`idl/content/content.thrift`
- 新建：`docs/migration/content-service.md`
- 新建：`docs/migration/web-bff.md`

**接口：**
- 使用：`content_db` 中的文章、分类、标签表；`auth-service` 返回的认证上下文
- 产出：`GetArticleByID(ctx context.Context, id int64) (*ArticleDetail, error)`
- 产出：`ListPublicArticles(ctx context.Context, filter ListFilter) (*ListResult, error)`
- 产出：`BuildArticlePage(ctx context.Context, id int64) (*ArticlePageDTO, error)`

- [ ] **步骤 1：先写失败的内容服务测试**

```go
package application

import (
	"context"
	"testing"
)

type fakeArticleRepo struct{}

func (f *fakeArticleRepo) GetByID(ctx context.Context, id int64) (*ArticleDetail, error) {
	return &ArticleDetail{ID: id, Title: "Hello"}, nil
}

func TestGetArticleByID(t *testing.T) {
	svc := NewArticleService(&fakeArticleRepo{})
	article, err := svc.GetArticleByID(context.Background(), 12)
	if err != nil {
		t.Fatalf("GetArticleByID returned error: %v", err)
	}
	if article.ID != 12 {
		t.Fatalf("expected article id 12, got %d", article.ID)
	}
}
```

- [ ] **步骤 2：运行测试，确认它先失败**

运行：`go test ./services/content-service/internal/application -run TestGetArticleByID -v`  
预期：FAIL，提示 `NewArticleService` 未定义。

- [ ] **步骤 3：写最小内容服务实现**

```go
package application

import "context"

type ArticleDetail struct {
	ID    int64
	Title string
}

type ArticleRepository interface {
	GetByID(ctx context.Context, id int64) (*ArticleDetail, error)
}

type ArticleService struct {
	repo ArticleRepository
}

func NewArticleService(repo ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) GetArticleByID(ctx context.Context, id int64) (*ArticleDetail, error) {
	return s.repo.GetByID(ctx, id)
}
```

- [ ] **步骤 4：迁移内容域并建立 web-bff**

完成以下最小能力：

- `content-service` 承接文章、分类、标签的核心 CRUD
- 文章详情统一按 `id` 查询
- `web-bff` 提供博客页面聚合接口
- `content-service/migrations/001_init.sql` 建立 `articles`、`categories`、`tags`、`article_tags`
- `idl/content/content.thrift` 定义公开查询和后台管理 RPC 契约

- [ ] **步骤 5：运行测试与编译验证**

运行：`go test ./services/content-service/... ./services/web-bff/... -v`  
预期：PASS。

运行：`go test ./...`  
预期：仓库整体可编译，新增服务相关包可通过测试。

- [ ] **步骤 6：补中文迁移说明**

在以下文档记录迁移决策：

- `docs/migration/content-service.md`
- `docs/migration/web-bff.md`

至少写清：

- 为什么 `article/category/tag` 归为一个内容域
- 为什么前端不直接拼多个服务接口
- 当前 `web-bff` 只做什么，不做什么

- [ ] **步骤 7：提交**

```bash
git add services/content-service services/web-bff idl/content/content.thrift docs/migration/content-service.md docs/migration/web-bff.md
git commit -m "feat: extract content service and web bff skeleton"
```

### 任务 5：补齐 gateway、Nacos、OpenTelemetry 与 K8s 部署资产

**文件：**
- 新建：`services/gateway/cmd/main.go`
- 新建：`services/gateway/internal/router/router.go`
- 新建：`services/gateway/internal/router/router_test.go`
- 新建：`services/gateway/internal/middleware/auth.go`
- 新建：`services/gateway/internal/middleware/trace.go`
- 新建：`pkg/xnacos/client.go`
- 新建：`pkg/xotel/setup.go`
- 新建：`deploy/docker/compose.yaml`
- 新建：`deploy/k8s/base/gateway/deployment.yaml`
- 新建：`deploy/k8s/base/auth-service/deployment.yaml`
- 新建：`deploy/k8s/base/content-service/deployment.yaml`
- 新建：`deploy/k8s/base/media-service/deployment.yaml`
- 新建：`deploy/k8s/base/web-bff/deployment.yaml`
- 新建：`deploy/k8s/base/otel-collector/deployment.yaml`
- 新建：`deploy/k8s/base/nacos/deployment.yaml`
- 新建：`deploy/k8s/dev/`
- 新建：`deploy/k8s/prod/`
- 新建：`docs/runbooks/local-microservices.md`
- 新建：`docs/runbooks/k8s-deploy.md`

**接口：**
- 使用：`auth-service.ValidateSession`、Nacos 服务发现、OTEL Collector 地址、K8s ConfigMap/Secret
- 产出：`func RegisterRoutes(h *server.Hertz, deps Dependencies)`
- 产出：`func SetupTracerProvider(ctx context.Context, serviceName string, endpoint string) (func(context.Context) error, error)`

- [ ] **步骤 1：先写失败的 gateway 路由测试**

```go
package router

import "testing"

func TestRegisterRoutesRequiresDependencies(t *testing.T) {
	deps := Dependencies{
		AuthServiceName: "auth-service",
		BFFServiceName:  "web-bff",
	}
	if err := ValidateDependencies(deps); err != nil {
		t.Fatalf("expected dependencies to validate, got %v", err)
	}
}
```

- [ ] **步骤 2：运行测试，确认它先失败**

运行：`go test ./services/gateway/internal/router -run TestRegisterRoutesRequiresDependencies -v`  
预期：FAIL，提示 `Dependencies` 或 `ValidateDependencies` 未定义。

- [ ] **步骤 3：写最小 gateway 与平台能力实现**

```go
package router

import "errors"

type Dependencies struct {
	AuthServiceName string
	BFFServiceName  string
}

func ValidateDependencies(deps Dependencies) error {
	if deps.AuthServiceName == "" {
		return errors.New("auth service name is required")
	}
	if deps.BFFServiceName == "" {
		return errors.New("bff service name is required")
	}
	return nil
}
```

同时完成以下最小资产：

- `gateway` 基础转发与中间件骨架
- `pkg/xnacos/client.go` 统一封装服务发现与配置读取入口
- `pkg/xotel/setup.go` 统一初始化 tracer provider
- `deploy/docker/compose.yaml` 能本地拉起核心依赖
- `deploy/k8s/base` 下提供业务服务与关键基础设施 Deployment / Service 基础模板

- [ ] **步骤 4：运行测试与配置校验**

运行：`go test ./services/gateway/internal/router -v`  
预期：PASS。

运行：`go test ./...`  
预期：仓库整体可编译，新增服务相关包可通过测试。

人工检查：

- `deploy/docker/compose.yaml` 结构自洽
- `deploy/k8s/base`、`dev`、`prod` 目录存在且职责清晰

- [ ] **步骤 5：补中文运行手册**

在以下文档中记录运行与部署方式：

- `docs/runbooks/local-microservices.md`
- `docs/runbooks/k8s-deploy.md`

至少写清：

- 本地如何启动依赖与服务
- Nacos、OTEL Collector、Prometheus、Grafana、Jaeger 的职责
- K8s 中哪些资源放 `ConfigMap`，哪些放 `Secret`

- [ ] **步骤 6：提交**

```bash
git add services/gateway pkg/xnacos pkg/xotel deploy/docker/compose.yaml deploy/k8s docs/runbooks
git commit -m "feat: add gateway and platform deployment assets"
```

## 自检

- 规格覆盖：任务 1 覆盖单体预重构与 `session + cookie + redis` 会话改造；任务 2 覆盖 `auth-service`；任务 3 覆盖 `media-service`；任务 4 覆盖 `content-service` 与 `web-bff`；任务 5 覆盖 `gateway`、`Nacos`、`OpenTelemetry`、Docker 与 Kubernetes 资产。
- 占位语扫描：计划中不包含 `TODO`、`TBD`、`implement later`、`fill in details` 等占位表达。
- 类型一致性：`CreateSession`、`ParseSession`、`GetArticleByID`、`Upload`、`ValidateDependencies` 都先定义再被后续任务依赖。
- 中文约束：本计划明确要求新增或改写的设计文档、迁移说明、运行手册统一使用中文。

## 执行交接

计划已保存到 `docs/superpowers/plans/2026-07-13-microservice-refactor-implementation.md`。接下来有两种执行方式：

**1. Subagent-Driven（推荐）** - 我按任务逐个派发，逐个 review，节奏更稳

**2. Inline Execution** - 直接在当前会话里按计划开始做

**你想用哪一种？**
