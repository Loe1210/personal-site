# Phase 11 - 阶段 A/B 首轮收口与 Swagger 清理

## Goal

围绕 `docs/pending-tasks.md` 中的阶段 A 与阶段 B，先完成一轮可交付收口：补回上下篇导航、移除 Swagger、补齐基础配置拆分、增加第一批回归测试，并提供 Docker 交付方式。

---

## 本阶段完成项

### 1. 恢复博客详情页上下篇导航

- 在后端 article service 中补回相邻文章查询能力
- 新增公开接口：`GET /api/articles/id/:id/adjacent`
- 前端 `BlogAPI.getAdjacentPosts` 不再返回空占位，而是改为真实请求后端接口
- 详情页导航区现在可以根据当前文章显示上一篇和下一篇

### 2. 清理 Swagger 遗留能力

- 移除 `main.go` 中 `/swagger.json` 与 `/swagger.yaml` 暴露
- 从路由注册中移除站点层 Swagger 入口
- 删除 `biz/site/handler.go` 中只服务 Swagger 的页面逻辑
- 删除 `docs/docs.go`、`docs/swagger.json`、`docs/swagger.yaml`
- 同时移除 `main.go` 顶部已失效的 Swagger 注释头

### 3. 完成第一轮黑线修复与后台视觉收口

- 前台博客页背景从固定附件改为普通滚动背景
- 收紧横向溢出与纵向滚动行为，降低滚动时出现黑线的概率
- 后台补齐 header actions、panel 容器层次、tab 阴影与布局细节
- 移除 admin 页面中重复加载的 `marked.js`

**说明**：这部分属于低风险样式修复，代码已落地，但最终视觉效果仍建议人工打开页面滚动复核。

### 4. 完成配置拆分

- `configs/config.go` 从原有 `server/mysql/session` 扩展为：
  - `server`
  - `mysql`
  - `session`
  - `upload`
  - `site`
- 保留原有环境变量覆盖方式，并新增：
  - `UPLOAD_ROOT_DIR`
  - `UPLOAD_PUBLIC_BASE_PATH`
  - `UPLOAD_MAX_IMAGE_SIZE_MB`
  - `SITE_TITLE`
  - `SITE_BASE_URL`
- `service/upload.go` 已接入新的 upload 配置，而不是继续完全写死上传目录与大小限制

### 5. 补齐第一批测试

新增测试覆盖：

- `service/article_adjacent_test.go`
  - 验证相邻文章查找逻辑
- `configs/config_test.go`
  - 验证默认值
  - 验证 YAML 加载
  - 验证环境变量覆盖

### 6. 补齐第一版 Docker 交付

新增文件：

- `Dockerfile`
- `docker-compose.yml`
- `README.md`

目标：提供最小可运行容器交付，不引入额外反向代理。

---

## 验证结果

本阶段已执行并通过：

```bash
go test ./...
go build ./...
```

### Docker 验证说明

尝试执行：

```bash
docker build -t personal-site .
```

但当前执行环境中没有可用的 Docker 命令，因此本阶段无法在本机完成镜像构建验证，后续需要在安装了 Docker 的环境中补一次实际构建检查。

---

## 当前仍保留的后续事项

1. RBAC 管理接口仍未进入实现，本阶段按计划继续延后
2. 滚动黑线问题与后台视觉对齐仍建议人工页面复核
3. 测试虽然补上了第一批关键路径，但整体覆盖率仍然偏低
4. Docker 文件已落地，但还缺一次真实镜像构建验证
