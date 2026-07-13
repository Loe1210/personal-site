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
- `UPLOAD_ROOT_DIR`
- `UPLOAD_PUBLIC_BASE_PATH`
- `UPLOAD_MAX_IMAGE_SIZE_MB`
- `SITE_TITLE`
- `SITE_BASE_URL`

## Docker 运行

构建镜像：

```bash
docker build -t personal-site .
```

使用 Compose 启动应用和 MySQL：

```bash
docker compose up --build
```

启动后访问：`http://localhost:8888`

## 验证

项目当前可通过以下命令做基础验证：

```bash
go test ./...
go build ./...
```
