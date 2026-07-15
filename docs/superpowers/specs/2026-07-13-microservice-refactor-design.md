# 微服务重构设计文档

## 目标

在保留现有个人站点核心功能可持续演进的前提下，将当前 Hertz 单体项目重构为一套以 `Hertz + Kitex + Nacos + OpenTelemetry + Kubernetes` 为核心技术栈的微服务系统。

本次设计目标不是“一次性重写全部业务”，而是定义一条可落地、可分阶段验收、适合学习完整微服务基础设施的演进路径。

## 本轮范围

本轮设计覆盖以下内容：

- 微服务拆分后的仓库目录结构
- `gateway`、`auth-service`、`content-service`、`media-service`、`web-bff` 的职责边界
- 对外 HTTP API 与内部 RPC 的边界划分
- `Nacos`、`OpenTelemetry`、`Kubernetes` 的落地方式
- 从当前单体迁移到微服务的五阶段实施顺序

本轮明确不做：

- 直接开始代码级全量重写
- 一开始就引入 Service Mesh、Operator、灰度发布体系
- 将 `article`、`category`、`tag` 继续细拆成更多小服务
- 为所有能力设计生产级高可用拓扑

## 当前项目背景

当前仓库是一个 Go 单体项目，技术栈以 `Hertz + GORM + MySQL` 为主，前台博客与后台管理功能已经可用。核心业务模块已经具备初步领域边界：

- `auth`
- `article`
- `category`
- `tag`
- `upload`
- `site`

当前代码存在以下微服务演进障碍：

1. 业务边界虽然存在，但服务层仍有跨域直连，例如文章逻辑直接查询标签和分类。
2. 认证体系当前只适合单体内使用，后续需要演进为 `session + cookie + redis` 共享会话方案，才能支撑多服务、多实例场景。
3. 中间件直接访问数据库，导致网关层和认证能力未来难以解耦。
4. 数据库初始化、迁移、种子数据写入职责混在一起，难以拆成按服务拥有数据的模式。
5. 当前目录结构更偏单体分层，不适合长期维护多服务并行演进。

## 备选方案

### 方案 1：激进微拆

将 `article`、`category`、`tag`、`upload`、`auth`、`rbac`、`site` 全部拆成独立服务。

优点：

- 形式上最像微服务
- 每个模块边界看起来最细

缺点：

- 当前项目体量偏小，内容域被过度拆分
- `article/category/tag` 会产生大量同步 RPC
- 联调、权限控制、数据一致性复杂度会陡增

### 方案 2：中等粒度拆分

拆成 `gateway`、`auth-service`、`content-service`、`media-service`、`web-bff` 五类服务。

优点：

- 既能练到完整微服务基础设施，也能控制边界复杂度
- 内容域仍然保持一致性
- 拆分顺序清晰，可渐进迁移

缺点：

- `web-bff` 和 `gateway` 的职责需要严格控制，避免膨胀

### 方案 3：保守演进

仅拆 `auth` 和 `upload`，其余继续保留在模块化单体中。

优点：

- 风险最低
- 迁移成本最小

缺点：

- 无法充分覆盖完整微服务基础设施的学习目标
- 核心内容域仍然停留在单体

## 采用方案

采用方案 2。

最终目标架构为：

- `gateway`
- `auth-service`
- `content-service`
- `media-service`
- `web-bff`

原因：

1. `auth` 是天然横切能力，必须独立。
2. `upload/media` 与内容主域耦合较低，适合早拆。
3. `article/category/tag` 当前属于同一个内容域，拆成一个 `content-service` 最合理。
4. `gateway + web-bff` 的组合既能承接入口治理，也能给前端提供稳定聚合接口。

## 设计说明

### 1. 仓库结构

仓库采用 `monorepo + 多服务并列` 组织方式：

```text
personal-site/
  services/
    gateway/
    auth-service/
    content-service/
    media-service/
    web-bff/

  idl/
    auth/
    content/
    media/

  pkg/
    xconfig/
    xlog/
    xtrace/
    xerror/
    xresponse/
    xotel/
    xnacos/

  deploy/
    docker/
    k8s/
      base/
      dev/
      prod/
    helm/

  docs/
    architecture/
    runbooks/
    migration/

  go.work
  Makefile
  README.md
```

每个服务内部统一为：

```text
services/<service-name>/
  cmd/
  internal/
    handler/
    application/
    domain/
    repository/
    infra/
    rpc/
    config/
  migrations/
  configs/
  Dockerfile
```

### 2. 服务职责

#### `gateway`

- 外部唯一 HTTP 入口
- 统一中间件：CORS、限流、request id、trace 注入、基础鉴权
- 路由转发到下游服务或 `web-bff`

#### `auth-service`

- 登录、登出、用户信息
- 登录、登出、session 创建、session 校验
- 角色与权限校验
- 向其他服务提供认证与授权 RPC

#### `content-service`

- 文章 CRUD
- 分类 CRUD
- 标签 CRUD
- 公开文章查询与后台内容管理
- 文章详情主键统一按 `id` 查询

#### `media-service`

- 文件上传
- 文件元数据存储
- 文件访问地址生成
- 对象存储适配层

#### `web-bff`

- 为前台博客页和后台管理页提供聚合接口
- 将多个下游服务返回值组装成前端稳定 DTO
- 不拥有业务真相，只做聚合和编排

### 3. API / RPC 边界

#### 对外 HTTP API

统一通过 `gateway` 暴露：

```text
POST   /api/auth/login
POST   /api/auth/logout
GET    /api/auth/me

GET    /api/articles
GET    /api/articles/{id}
GET    /api/categories
GET    /api/tags

GET    /api/admin/articles
POST   /api/admin/articles
PUT    /api/admin/articles/{id}
DELETE /api/admin/articles/{id}

POST   /api/admin/upload
```

#### `auth-service` RPC

```text
CreateSession(username, password) -> session_id / cookie / user
ValidateSession(session_id) -> user_id / username / roles / permissions
RefreshSession(session_id) -> renewed session metadata
CheckPermission(user_id, permission_code) -> allowed
GetUser(user_id) -> user profile
```

#### `content-service` RPC

```text
ListPublicArticles(filter) -> article list
GetArticleByID(id) -> article detail
ListAdminArticles(filter) -> article list
CreateArticle(dto) -> article
UpdateArticle(dto) -> article
DeleteArticle(id) -> result

ListCategories() -> category list
CreateCategory(dto) -> category
UpdateCategory(dto) -> category
DeleteCategory(id) -> result

ListTags() -> tag list
CreateTag(dto) -> tag
UpdateTag(dto) -> tag
DeleteTag(id) -> result
```

#### `media-service` RPC

```text
UploadFile(meta, stream_ref) -> file_id / url
GetFile(file_id) -> metadata
DeleteFile(file_id) -> result
```

### 4. 数据边界

数据库按服务拥有，不将 MySQL 视为业务服务目录。

- `auth-service` 拥有 `auth_db`
- `content-service` 拥有 `content_db`
- `media-service` 拥有 `media_db`

数据库部署资源位于 `deploy/`，迁移脚本位于各服务自己的 `migrations/`。

禁止跨服务直接读取彼此数据库表。

### 5. Nacos 落地

`Nacos` 同时承担：

- 服务注册与发现
- 配置中心

所有服务启动后向 `Nacos` 注册服务名：

- `gateway`
- `auth-service`
- `content-service`
- `media-service`
- `web-bff`

配置按服务维度管理，例如：

- `gateway-dev.yaml`
- `auth-service-dev.yaml`
- `content-service-dev.yaml`
- `media-service-dev.yaml`
- `web-bff-dev.yaml`

按环境区分 `dev / test / prod`。

### 6. OpenTelemetry 落地

可观测性优先落 `trace + metrics`，链路为：

```text
Service -> OpenTelemetry SDK -> OTEL Collector -> Jaeger/Tempo + Prometheus -> Grafana
```

埋点重点：

- `gateway`：HTTP span、状态码、耗时、request id、session 校验调用
- `auth-service`：登录、session 校验、权限查询、Redis 会话访问、数据库访问
- `content-service`：文章查询/写入、分类标签查询、数据库访问
- `media-service`：上传、对象存储操作、数据库访问
- `web-bff`：聚合请求与下游 RPC 调用

日志第一版统一要求：

- 输出 `trace_id`
- 输出 `service name`
- 错误日志包含 `request_id` 与必要上下文

### 7. Kubernetes 落地

部署结构采用：

```text
deploy/k8s/
  base/
    gateway/
    auth-service/
    content-service/
    media-service/
    web-bff/
    nacos/
    mysql-auth/
    mysql-content/
    mysql-media/
    redis/
    minio/
    otel-collector/
    prometheus/
    grafana/
    jaeger/
  dev/
  prod/
```

原则：

- 外部流量只进入 `gateway`
- 下游业务服务全部采用集群内部 `ClusterIP`
- 配置使用 `ConfigMap`
- 密钥使用 `Secret`
- 每个服务提供健康检查与优雅退出能力

### 8. 实施阶段

#### 阶段一：单体预重构

- 改造目录边界
- 将单体内会话能力重构为 `session + cookie + redis` 共享会话方案
- 拆分数据库初始化职责
- 抽取公共基础设施包

#### 阶段二：拆出 `auth-service` 与 `media-service`

- 登录、权限、上传独立为服务
- 引入第一批 `Kitex` RPC
- 建立 `auth_db` 与 `media_db`

#### 阶段三：拆出 `content-service`

- 迁移文章、分类、标签等内容能力
- 建立 `content_db`
- 初版建立 `web-bff`

#### 阶段四：补齐治理与可观测性

- 引入 `gateway`
- 接入 `Nacos`
- 接入 `OpenTelemetry + OTEL Collector + Jaeger/Tempo + Prometheus + Grafana`

#### 阶段五：K8s 化与工程化交付

- 所有服务完成容器化
- 部署基础设施与业务服务到 K8s
- 建立 `dev / prod` 环境差异化部署

## 错误处理

- 认证失败由 `auth-service` 输出明确错误码，`gateway` 只负责透传或做基础映射。
- 网关层只负责读取 cookie 中的 session 标识并调用 `auth-service` 校验，不在网关中直接查数据库或直查 Redis。
- 下游服务不可用时，`web-bff` 与 `gateway` 需要返回统一错误响应并记录 trace。
- 各服务需要提供 `/health` 或等价健康检查接口，供 K8s readiness / liveness 使用。

## 验证方式

- 阶段一以单体回归测试和 `session + cookie + redis` 会话改造验证为主
- 阶段二开始验证 `auth-service`、`media-service` 的独立可运行性与 RPC 调用
- 阶段三验证 `content-service` 对文章主链路的完整承接
- 阶段四验证 `Nacos` 服务发现、配置加载、Jaeger/Tempo 链路追踪、Grafana 指标面板
- 阶段五验证 K8s 部署、滚动更新、健康检查和基础扩容行为

## 完成标准

- 仓库完成从单体目录到多服务目录的重构
- 认证、内容、媒体能力均拥有独立服务边界与独立数据库
- 对外请求统一经 `gateway`
- 服务间调用通过 `Kitex + Nacos`
- 全链路可通过 `OpenTelemetry` 追踪
- 系统可在 Kubernetes 环境中部署运行
