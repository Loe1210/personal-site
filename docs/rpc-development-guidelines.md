# RPC 开发规范

本文是当前项目的 RPC 主线规范。当前项目统一使用 `proto + Kitex`，不再把 thrift 作为主线接口定义方式。

## 基本原则

- RPC 接口必须先定义 IDL，再生成代码，再实现业务逻辑。
- IDL 文件统一放在仓库 `idl/` 目录下，按服务拆分，例如 `idl/content/content.proto`。
- 生成代码只读，禁止手动修改 `kitex_gen/` 下的生成文件。
- HTTP 请求和页面接口继续手写 Hertz 路由；RPC 只负责服务间契约，不替代所有 HTTP handler。
- 没有真实跨服务调用需求时，只保留 RPC server 能力，不强行写 client。

## Proto 字段规则

- 字段序号一旦发布，不允许随意修改。
- 删除字段时不要复用旧字段序号，避免二进制序列化/反序列化数据错位。
- 新增字段必须使用新的字段序号，并考虑兼容老客户端。
- message 命名表达业务语义，不要把前端页面 view model 全部搬进 proto。

## 生成代码规则

- 通过 Makefile 或脚本统一生成 Kitex 代码。
- 生成文件变更必须来自 IDL 变更，不允许手写补丁。
- 代码评审时优先看 `.proto` 的接口语义，再看生成结果。
- 如果生成结果不符合预期，应修改 `.proto` 或生成脚本，而不是修改 `kitex_gen/`。

## 当前项目边界

- `auth-service`：负责登录、会话、认证相关能力。
- `content-service`：负责文章、分类、标签内容能力。
- `media-service`：负责媒体上传和文件元数据能力。
- `web-bff`：负责面向前端的聚合能力，当前不提前扩大职责。
- `gateway`：负责外部 HTTP 入口与路由转发，不承载业务逻辑。

## 与 HTTP DTO 的关系

- HTTP 请求/响应 DTO 可以继续由 Hertz handler 手写维护。
- RPC DTO 由 proto 生成。
- 不要为了“统一”而把所有 HTTP DTO 都改成 proto 生成对象；这会让页面展示模型和服务契约过早耦合。

## 当前执行策略

1. 先保证 Docker Compose 下前端、网关、服务、数据库、Redis 能稳定运行。
2. 保持 RPC server 骨架可编译、可启动。
3. 等出现明确的服务间调用场景，再补 RPC client。
4. 每次修改 IDL 后，必须重新生成代码并运行 `go test ./...`。