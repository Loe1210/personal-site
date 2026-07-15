# 个人小站

这是一个基于 Go + Hertz + MySQL 的个人小站项目。当前后端已经从单体入口迁移到微服务运行方式，外部 HTTP 流量统一从 `gateway` 进入，再转发到 `auth-service`、`content-service`、`media-service` 和 `web-bff`。

## 当前架构

- `services/gateway`：统一 HTTP 入口，默认监听 `http://127.0.0.1:8888`
- `services/auth-service`：登录、登出、当前用户、权限校验基础能力
- `services/content-service`：文章内容域，包含文章列表和文章详情等接口
- `services/media-service`：上传和文件元数据能力
- `services/web-bff`：面向前台页面的聚合层
- `frontend`：Nginx 静态前端入口，默认访问 `http://127.0.0.1:8080`
- `deploy/docker/compose.yaml`：本地微服务运行环境

## 本地运行

推荐使用微服务版 Docker Compose：

```bash
make micro-up
```

如果本机没有 `make`，可以执行等价命令：

```bash
docker compose -f deploy/docker/compose.yaml up -d --build
```

启动后优先访问前端页面：

```text
http://127.0.0.1:8080
```

后端 API 网关仍然保留在：

```text
http://127.0.0.1:8888
```

## 验证

基础 Go 测试：

```bash
go test ./...
```

端到端 smoke 验证：

```bash
make micro-smoke
```

如果本机没有 `make`，可以执行：

```powershell
powershell -ExecutionPolicy Bypass -File scripts/smoke/microservices_smoke.ps1
```

smoke 会验证 gateway 健康检查、auth 未登录态、content 文章列表，以及 `session cookie + Redis` 登录态闭环。

## 停止本地环境

```bash
make micro-down
```

等价命令：

```bash
docker compose -f deploy/docker/compose.yaml down
```

## 配置说明

配置支持通过 YAML 和环境变量覆盖。常用环境变量包括：

- `APP_HOST`
- `APP_PORT`
- `MYSQL_HOST`
- `MYSQL_PORT`
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_DBNAME`
- `MYSQL_CHARSET`
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

## 认证说明

认证使用 `session cookie + Redis`：

- 登录成功后服务端写入 session cookie
- session 内容存储在 Redis 中
- 受保护接口通过 cookie 中的 session id 查询 Redis 并恢复用户信息
- 默认本地 smoke 使用 `admin/admin` 验证登录链路

## 遗留单体说明

旧单体入口、`biz/`、`service/`、`dal/db/` 和根目录旧 IDL 已删除。后续新增功能应优先落在对应 `services/*` 服务内，不再恢复旧单体分层。
