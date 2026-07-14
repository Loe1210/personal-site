# 个人小站

这是一个基于 Go + Hertz + MySQL 的个人小站项目，当前前台与后台都采用静态 SPA 方式交付，后端主要负责 API、认证、上传与静态资源分发。

## 本地直接运行

1. 按需修改 `configs/config.yaml`
2. 启动 MySQL，并确保数据库存在
3. 运行：

```bash
go run .
```

默认监听地址：`http://localhost:8888`

## 配置说明

当前配置按领域拆分为：

- `server`
- `mysql`
- `session`
- `session_store`
- `redis`
- `upload`
- `site`

同时支持通过环境变量覆盖关键配置，例如：

- `APP_HOST`
- `APP_PORT`
- `MYSQL_HOST`
- `MYSQL_PORT`
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_DBNAME`
- `MYSQL_CHARSET`
- `SESSION_SECRET`
- `SESSION_STORE_PREFIX`
- `SESSION_STORE_EXPIRE_HOUR`
- `SESSION_STORE_COOKIE_NAME`
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `REDIS_DB`
- `UPLOAD_ROOT_DIR`
- `UPLOAD_PUBLIC_BASE_PATH`
- `UPLOAD_MAX_IMAGE_SIZE_MB`
- `SITE_TITLE`
- `SITE_BASE_URL`

## Docker 运行

当前容器化运行入口已经切换为微服务版 Compose：

```bash
make micro-up
```

该命令会使用 `deploy/docker/compose.yaml` 启动 MySQL、Redis、Nacos、OTel Collector 以及各业务服务。

如果本机没有 `make`，可以执行等价命令：

```bash
docker compose -f deploy/docker/compose.yaml up -d --build
```

启动后优先访问 gateway：`http://localhost:8888`

运行 smoke 验证：

```bash
make micro-smoke
```

关闭本地微服务环境：

```bash
make micro-down
```

## 验证

项目当前可通过以下命令做基础验证：

```bash
go test ./...
go build ./...
```


## 认证说明

当前正在为微服务拆分做认证链路预重构，现阶段使用 `session + cookie + redis` 作为目标会话模型。

- 登录成功后会返回会话信息，并写入浏览器 cookie
- 调试受保护接口时，需要携带浏览器返回的 session cookie
- 当前仓库仍处于单体阶段，Redis 共享会话的完整接入会在后续 `auth-service` 拆分阶段完成
