# 当前待办与后续开发计划

本文只记录当前微服务版本后续要做的事情，旧单体计划不再作为执行依据。

## 当前已完成并可继续演进的基础

- 本地入口已经切到 `frontend` 容器，浏览器访问 `http://127.0.0.1:8080`。
- 前端静态资源由 Nginx 承载，`/api` 请求由 Nginx 转发到 `gateway`。
- `gateway` 对外 HTTP 入口为 `http://127.0.0.1:8888`。
- 后端服务拆分为 `auth-service`、`content-service`、`media-service`、`web-bff`。
- 登录鉴权采用 `session cookie + Redis`，不使用 JWT 作为当前主线方案。
- `content-service` 已恢复文章、分类、标签相关接口，分类和标签模型已从文章模型中拆出。
- RPC 当前以 `proto + Kitex` 为规范方向，服务端骨架已具备；没有真实跨服务调用需求前，不强行引入 RPC client。

## 近期优先级

1. 完成前端页面手动联调：前台文章列表、分类、标签、详情页；后台登录、文章、分类、标签管理。
2. 给 `content-service` 的分类和标签接口补充更明确的自动化测试，避免之后再次漏掉接口。
3. 梳理后台文章创建/编辑/删除流程，确认前端请求路径与网关、内容服务完全一致。
4. 确认上传链路是否继续走本地存储，还是进入 MinIO/对象存储方案。
5. 等出现真实服务间协作场景后，再补 Kitex RPC client，不提前制造耦合。

## 中期任务

- 补齐 RBAC 管理接口和后台权限管理页面。
- 为每个服务继续收敛自己的 `internal/model`、`internal/service`、`internal/dal` 边界。
- 将旧单体遗留代码逐步清理干净，只保留当前微服务主线需要的目录。
- 完善 Nacos、OpenTelemetry、Prometheus、Grafana 的本地观测链路说明。
- 准备 K8s 部署清单，但不要在本地 Docker Compose 稳定前提前切换主线。

## 日常验证命令

```powershell
go test ./...
C:\ProgramData\chocolatey\bin\make.exe micro-smoke
```

## 当前访问地址

- 前端页面：`http://127.0.0.1:8080`
- 网关接口：`http://127.0.0.1:8888`
- Nacos：`http://127.0.0.1:8848/nacos`
- Jaeger：`http://127.0.0.1:16686`
- Prometheus：`http://127.0.0.1:9090`
- Grafana：`http://127.0.0.1:3000`