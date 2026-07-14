# Proto + Kitex 微服务目录规范设计

## 目标

将当前个人站点的微服务代码规范统一为 `proto + Kitex` 的契约体系，并把每个服务的目录边界固定为：`biz / cmd / internal / pkg / idl`。HTTP 继续手写，RPC 继续走 Kitex；`proto` 只负责服务契约与消息结构，不承载业务实现。

## 本轮范围

本轮设计只覆盖下面这些内容：

- `idl` 从 `thrift` 全量切换到 `proto`
- 每个服务的目录结构统一
- 每个服务保留独立的 `pkg/`
- 根目录不再保留公共 `pkg/`
- `proto` 继续给 `Kitex` 使用
- 生成代码只读，禁止手改
- RPC 与 HTTP 的职责边界
- `Makefile` 的 proto 生成与验证方式
- 分阶段迁移顺序

本轮不做：

- 直接开始大规模代码重写
- 直接把所有服务拆到更多子服务
- 直接切换到 gRPC
- 直接引入 Service Mesh 或 Operator

## 当前项目背景

项目已经完成从单体向微服务的架构迁移，服务拆分已经存在，后续重点是把契约层和目录层统一起来。当前主要问题是：

1. `thrift` 仍然是当前契约主线，但后续要切到 `proto`。
2. 各服务的工具代码边界还需要进一步收紧，避免形成跨服务耦合。
3. 目录结构需要固定成一致的服务模板，方便后续继续扩展。
4. RPC 和 HTTP 的责任边界需要进一步明确，防止生成代码和手写代码混在一起。

## 采用方案

### 方案 1：继续使用 thrift

优点：
- 与当前已有代码一致
- 迁移成本最低

缺点：
- 不符合当前已确认的 proto 方向
- 后续继续学习和对接 Kitex 时还要再迁一次

### 方案 2：proto + Kitex，按服务统一目录

优点：
- 与最终服务契约方向一致
- 后续服务间 RPC 可以直接延续 Kitex
- 目录边界清晰，便于长期维护

缺点：
- 需要一次性整理 IDL、生成命令和目录结构

### 方案 3：proto + gRPC

优点：
- 生态成熟

缺点：
- 与当前项目既定的 Kitex 路线不一致
- 会增加额外迁移成本

## 采用方案

采用方案 2。

最终约定：

- `proto` 继续给 `Kitex` 用
- `HTTP` 仍然由 Hertz 手写
- 每个服务独立拥有自己的 `pkg/`
- 根目录不再保留统一 `pkg/`
- 生成代码只读，禁止手改

## 目录规范

### 仓库总结构

```text
personal-site/
  services/
    auth-service/
    content-service/
    media-service/
    gateway/
    web-bff/

  idl/
    auth/
    content/
    media/

  deploy/
  docs/
  Makefile
  README.md
```

### 每个服务统一结构

```text
services/<service-name>/
  biz/
    handler/
    model/
    router/
    mw/
  cmd/
    configs/
    main.go
    router.go
  idl/
  internal/
    dal/
    model/
    service/
  pkg/
  migrations/
```

### 目录职责

- `biz/`：HTTP 入口层、路由、handler、HTTP DTO
- `cmd/`：程序入口、配置装配、启动 wiring
- `idl/`：当前服务的 RPC 契约定义
- `internal/model/`：服务内部共享对象
- `internal/service/`：业务逻辑
- `internal/dal/`：数据库、Redis、RPC client、存储适配
- `pkg/`：仅当前服务可使用的工具代码

## 关键约束

### 1. `pkg` 约束

- 不保留根目录 `pkg/`
- 每个服务保留自己的 `pkg/`
- `services/<service-name>/pkg/` 只能被本服务引用
- 不允许服务之间互相 import 对方的 `pkg/`

### 2. RPC 约束

- RPC 接口统一写在 `idl/*.proto`
- 生成代码只读，不手工修改生成文件
- 字段序号一旦确定，不得随意变更
- `proto` 生成代码只表达契约，不表达业务实现

### 3. HTTP 约束

- HTTP 路由继续手写
- `biz` 只负责 HTTP 入口，不放核心业务逻辑
- `biz` 的请求/响应对象和 RPC 契约可以相互映射，但不能混成一层

### 4. 服务边界约束

- `internal/model` 是服务内共享对象边界
- `internal/service` 只做业务编排
- `internal/dal` 只做基础设施访问
- 禁止跨服务直接访问数据库表

## IDL 设计

### 目录

```text
idl/
  auth/
  content/
  media/
```

### 服务映射

- `auth-service` 对应 `idl/auth/*.proto`
- `content-service` 对应 `idl/content/*.proto`
- `media-service` 对应 `idl/media/*.proto`

### 约定

- 每个 proto 文件独立 `package`
- 每个 proto 文件独立 `go_package`
- 同一服务内的消息结构优先复用，不重复定义
- 公共消息只有在稳定后才抽到 `common.proto`

## 生成目录设计

推荐每个 proto 对应独立生成包，避免互相覆盖：

```text
services/<service-name>/biz/model/<domain>/
```

示例：

```text
services/content-service/biz/model/article/
services/content-service/biz/model/auth/
services/media-service/biz/model/upload/
```

## 服务边界

### auth-service

职责：

- 登录
- 登出
- 当前用户
- session 校验
- 权限校验

认证方案：

- `session + cookie + redis`
- 不回退 JWT 作为主登录态

### content-service

职责：

- 文章 CRUD
- 分类 CRUD
- 标签 CRUD
- 公开文章查询
- 后台内容管理

### media-service

职责：

- 文件上传
- 文件元数据存储
- 文件访问地址生成

### gateway

职责：

- 对外统一入口
- 路由转发
- 基础鉴权与 trace 注入

### web-bff

职责：

- 聚合前端所需数据
- 不拥有业务真相

## RPC 命名建议

### auth.proto

- `UserLogin`
- `GetCurrentUser`
- `ValidateSession`
- `CheckPermission`

### content.proto

- `CreateArticle`
- `UpdateArticle`
- `DeleteArticle`
- `GetArticleByID`
- `GetArticleBySlug`
- `ListArticles`
- `ListAdminArticles`
- `CreateCategory`
- `ListCategories`
- `CreateTag`
- `ListTags`

### media.proto

- `UploadImage`
- `GetUploadInfo`
- `DeleteUpload`

## Makefile 约束

`Makefile` 需要补齐：

- proto 生成命令
- 生成物检查命令
- 全量测试命令
- smoke 验证命令

要求：

- 生成命令清晰可重复
- 生成物不可手改
- Proto 变更后可以一条命令重新生成

## 实施顺序

### 阶段 1：IDL 切换

- `thrift` 文件改成 `proto`
- 先定 `auth / content / media` 三个核心服务的契约
- 补齐 `go_package`、`package` 和生成目录

### 阶段 2：目录统一

- 统一各服务为 `biz / cmd / internal / pkg / idl`
- 每个服务的 `pkg` 独立存在
- 删除根 `pkg`

### 阶段 3：生成与接线

- 接入 Proto 生成命令
- 接入 Kitex 生成流程
- 保证 HTTP 手写与 RPC 生成物分离

### 阶段 4：清理遗留文档与旧契约

- 删除或标记过时的 thrift 说明
- 更新迁移文档和待办文档
- 统一 README 中的契约说明

## 验证方式

- `go test ./...`
- `make micro-smoke`
- proto 生成命令重复执行无冲突
- 生成目录不会出现交叉覆盖

## 完成标准

- 只保留 `proto` 作为 RPC 契约来源
- `Kitex` 继续作为 RPC 运行时
- 每个服务目录一致且边界清晰
- 每个服务独立 `pkg/`，不发生跨服务耦合
- 旧 `thrift` 路线不再作为主线