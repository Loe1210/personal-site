# 本地微服务运行手册

## 一键启动

在项目根目录执行：

```bash
make micro-up
```

该命令会通过 `deploy/docker/compose.yaml` 拉起 MySQL、Redis、Nacos、OTEL Collector、Prometheus、Grafana、Jaeger，以及 `auth-service`、`media-service`、`content-service`、`web-bff`、`gateway`。

## 数据库初始化

MySQL 容器启动时会执行：

```text
deploy/docker/init/001_databases.sql
```

它会创建三个服务数据库：

- `auth_db`
- `media_db`
- `content_db`

每个业务服务只能访问自己的 schema，不能跨服务直接读表。

## 常用命令

```bash
make micro-up
make micro-down
make micro-test
make micro-smoke
```

`make micro-smoke` 会执行 `scripts/smoke/microservices_smoke.ps1`，验证 gateway 健康检查、auth 登录 cookie 流程、content 文章列表等基础链路。

## 服务端口

- gateway: `http://127.0.0.1:8888`
- auth-service: `http://127.0.0.1:9001`
- media-service: `http://127.0.0.1:9002`
- content-service: `http://127.0.0.1:9003`
- web-bff: `http://127.0.0.1:9004`
- Nacos: `http://127.0.0.1:8848`
- OTEL Collector: `4317` / `4318`
- Prometheus: `http://127.0.0.1:9090`
- Grafana: `http://127.0.0.1:3000`
- Jaeger: `http://127.0.0.1:16686`

## 组件职责

- Nacos：服务注册发现和后续配置中心入口。
- OTEL Collector：统一接收各服务 OpenTelemetry 数据，再转发给追踪和指标系统。
- Prometheus：采集指标，用于服务健康和性能观测。
- Grafana：展示指标仪表盘。
- Jaeger：查看分布式链路追踪。
- Redis：保存 `session cookie` 对应的共享会话。
- MySQL：为每个服务提供独立 schema。