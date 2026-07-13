# Kubernetes 部署手册

## 目录结构

- `deploy/k8s/base`：所有环境共享的 Deployment、Service、ConfigMap、Secret 模板。
- `deploy/k8s/dev`：开发环境 overlay，默认单副本。
- `deploy/k8s/prod`：生产环境 overlay，gateway 默认双副本。

## 部署命令

开发环境：

```bash
kubectl apply -k deploy/k8s/dev
```

生产环境：

```bash
kubectl apply -k deploy/k8s/prod
```

## ConfigMap 和 Secret 分工

适合放入 ConfigMap 的内容：

- 服务名、服务端口、Nacos 地址。
- OTEL Collector 地址。
- Redis 地址。
- 非敏感的开关和公共配置。

必须放入 Secret 的内容：

- MySQL 密码。
- Redis 密码。
- session secret。
- 第三方对象存储密钥。
- 任何不能进入 Git 的凭据。

## 服务流量路径

外部 HTTP 流量统一进入 `gateway`，再转发或聚合到后端服务。前端页面接口优先访问 `web-bff`，领域数据由 `content-service`、`media-service`、`auth-service` 各自负责。

## 可观测性路径

业务服务通过 OpenTelemetry SDK 上报 trace 和 metrics 到 `otel-collector`。Collector 再把链路数据发往 Jaeger 或 Tempo，把指标交给 Prometheus，Grafana 负责展示。

## Nacos 落地方式

服务启动后通过统一封装 `pkg/xnacos` 注册自身实例，并通过 Nacos 查询下游服务地址。当前仓库先固定封装接口，后续接入真实 Nacos SDK 时，业务服务只需要继续调用封装层。