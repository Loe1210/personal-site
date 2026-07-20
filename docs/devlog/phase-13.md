# Phase 13 - 微服务本地运行、前端容器化与文档减噪

## Goal

本阶段目标是把前一轮微服务拆分成果推进到“本地可以跑、浏览器可以访问、后续开发不会被旧文档误导”的状态。重点包括：Docker Compose 本地运行链路、前端独立容器和 Nginx 转发、缺失的分类/标签接口恢复、proto + Kitex 规范收口，以及顶层开发文档清理。

---

## 本阶段完成项

### 1. 本地微服务运行链路

- Docker Desktop 与 WSL 环境已完成可用性排障。
- `deploy/docker/compose.yaml` 成为当前微服务本地启动入口。
- `make micro-up` 用于启动本地微服务栈。
- `make micro-smoke` 用于做基础健康检查。
- 浏览器主入口明确为 `http://127.0.0.1:8080`。
- 网关接口入口明确为 `http://127.0.0.1:8888`。

### 2. 前端独立容器与 Nginx 转发

- 新增 `frontend/` 容器构建目录。
- 前端静态页面由 Nginx 承载。
- Nginx 将 `/api` 请求转发到 `gateway`，避免浏览器直接感知内部服务地址。
- README 和本地运行 runbook 已更新当前访问方式。

### 3. 分类与标签接口恢复

- 恢复 `content-service` 的分类和标签接口，修复前台/后台访问 `/api/categories`、`/api/tags` 返回 404 的问题。
- 将 `Tag`、`Category` 从文章模型中拆开，分别放入内容服务自己的 `internal/model`。
- 补齐分类和标签的 service、repository、handler、router 链路。
- 前端 Nginx 增加兼容转发：`/api/categories`、`/api/tags` 转到网关的内容服务路径。

### 4. RPC 规范收口

- 当前项目主线确定为 `proto + Kitex`。
- thrift 不再作为当前主线 IDL。
- `kitex_gen/` 作为生成代码，只读维护。
- 当前阶段保留 RPC server 能力；没有真实跨服务调用需求前，不强行补 RPC client。
- HTTP 路由继续手写，RPC 用于后续服务间契约治理。

### 5. 文档清理

- 顶层旧单体计划和临时排障文档已清理，避免影响后续判断。
- 保留当前仍有价值的项目规范、RPC 规范、UI 规格、RBAC 规划、运行手册和历史开发日志。
- 新增本阶段开发日志，只记录本轮新增内容，不改动旧 phase 日志。

---

## 验证结果

本阶段已执行并通过：

```powershell
go test ./...
C:\ProgramData\chocolatey\bin\make.exe micro-smoke
```

本地接口曾验证通过：

```powershell
curl http://127.0.0.1:8080/api/categories
curl http://127.0.0.1:8080/api/tags
curl http://127.0.0.1:8080/api/admin/categories
curl http://127.0.0.1:8080/api/admin/tags
```

---

## 当前仍保留的后续事项

1. 浏览器内完整手动验证前台文章列表、文章详情、分类、标签和后台管理页面。
2. 为分类和标签接口补自动化测试，防止后续重构再次遗漏。
3. 等出现真实服务间调用需求后，再补 Kitex RPC client。
4. 继续推进 RBAC 管理接口、上传存储方案、K8s 部署清单和观测链路完善。
5. 继续清理旧单体遗留代码，但每次清理前必须确认不影响当前微服务运行链路。
