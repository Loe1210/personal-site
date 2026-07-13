# 本地微服务运行手册

## 启动基础依赖

在项目根目录执行：

```bash
docker compose -f deploy/docker/compose.yaml up -d
```

这会拉起 MySQL、Redis、Nacos、OTEL Collector、Prometheus、Grafana 和 Jaeger。业务服务仍然可以用 `go run` 分别启动，便于本地调试。

## 业务服务启动顺序

建议按以下顺序启动：

1. `auth-service`：负责登录、会话校验和权限判断。
2. `media-service`：负责上传文件和文件元数据。
3. `content-service`：负责文章、分类、标签内容域。
4. `web-bff`：负责页面聚合 DTO。
5. `gateway`：统一接收外部 HTTP 流量。

示例：

```bash
go run ./services/auth-service/cmd -config services/auth-service/configs/config.yaml
go run ./services/media-service/cmd -config services/media-service/configs/config.yaml
go run ./services/content-service/cmd -config services/content-service/configs/config.yaml
go run ./services/web-bff/cmd -config services/web-bff/configs/config.yaml
go run ./services/gateway/cmd -config configs/config.yaml
```

## 组件职责

- Nacos：服务注册发现和后续配置中心入口。
- OTEL Collector：统一接收各服务 OpenTelemetry 数据，再转发给追踪和指标系统。
- Prometheus：采集指标，用于服务健康和性能观测。
- Grafana：展示指标仪表盘。
- Jaeger：查看分布式链路追踪。
- Redis：保存 `session cookie` 对应的共享会话。
- MySQL：每个服务使用自己的 schema，禁止跨服务直接读表。

## 当前阶段说明

当前 Go 代码中的 `pkg/xnacos` 和 `pkg/xotel` 已经提供统一入口，但还没有接入真实 SDK exporter。这样做是为了先固定调用边界，避免在服务拆分阶段引入不可控依赖；后续可以在不改业务调用方的情况下替换为真实 Nacos 客户端和 OpenTelemetry exporter。