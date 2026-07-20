# Phase 14 - 博客界面刷新、宠物运行时与上传链路加固

## 日期

2026-07-20

## 对应提交

- 分支：`codex/blog-upload-pet-deploy-20260720`
- 提交：`97463f2 feat: refresh blog ui and harden uploads`

## 本轮真实更新范围

这轮更新不是单一 UI 调整，而是一次博客前台、后台上传、媒体服务、部署文档的组合更新。以下内容均来自本轮提交和本地验证结果，没有包含未实现或未上线的假设功能。

## 前台博客界面

- 重做博客首页布局，保留原有偏暗、山景背景、蓝色描边和玻璃质感风格，但改成更接近参考图的纵向文章流。
- 主页不再直接展示分类和标签侧栏，新增独立的分类页与标签页：`/blog/categories`、`/blog/tags`。
- 新增公共侧边栏脚本 `static/blog/js/sidebar.js`，让首页、分类页、标签页、文章页复用同一套导航入口。
- 更新文章列表脚本 `static/blog/js/list.js`，支持更紧凑的文章卡片、页面渐入、hover 反馈、分页和阅读全文跳转。
- 更新博客样式 `static/blog/css/blog.css`，加入标题区、文章卡片、导航、标签/分类页面、宠物层级、交互动效等样式。
- 新增可爱鼠标光标样式 `static/assets/css/cute-cursor.css`，并在博客/后台相关页面引入。

## 分类与标签

- 新增 `static/blog/categories.html` 和 `static/blog/js/categories.js`，分类页按分类聚合展示文章数量与文章列表。
- 新增 `static/blog/tags.html` 和 `static/blog/js/tags.js`，标签页保留“标签堆叠/掉落感”的展示方向。
- 前端 API 层 `static/blog/js/api.js` 补充了相邻文章接口调用能力，供文章详情页使用。

## 全局宠物

- 新增 `pet/` 静态资源目录，当前包含两只宠物：`yuexinmiao`、`zhangfei-tusun`。
- 新增 `pet/index.json`，用于列出可切换宠物。
- 新增 `static/blog/js/pet.js`，负责在页面中加载宠物、保存当前宠物选择、处理双击切换、拖动、空闲动作和全局显示。
- 前端镜像构建已把 `pet/` 复制到 nginx 静态目录，Go 服务镜像也同步复制 `pet/` 到 `/app/static/pet`。

## 后台管理与封面上传

- 后台文章编辑弹窗新增“从素材库选择”入口，同时保留本地图片上传。
- 本地封面上传改为调用后端分片上传接口，而不是旧的普通 `/api/media/upload` 接口。
- 前端在上传前使用浏览器 `crypto.subtle.digest('SHA-256')` 计算整文件 SHA-256，并在初始化上传任务时传给后端。
- 后端 `complete` 阶段仍会做最终文件 hash 校验，本轮没有取消校验。
- `static/admin/index.html` 已将后台脚本版本更新到 `/admin/js/admin.js?v=13`，避免浏览器继续使用旧缓存。

## 媒体服务修复

- 修复分片上传时 chunk 内容可能被读成空的问题：`UploadChunk` 改为读取 Hertz 缓冲后的 `c.Body()`，再用 `bytes.NewReader(body)` 交给 chunk 服务。
- 新增 `services/media-service/biz/upload/handler_test.go`，验证原始二进制 chunk 能被正确读取、保存大小正确、chunk SHA-256 正确。
- 修复分片合并后的文件权限：合并文件落盘后设置为 `0644`，避免 nginx 访问上传图片时出现 `403 Forbidden`。
- 修复缩略图生成后的文件权限：缩略图 JPEG 写入成功并关闭文件后设置为 `0644`。
- 更新 `merge_service_test.go`，验证合并后的文件至少具备 nginx 可读权限。

## 内容服务与文章详情

- 新增公开文章相邻篇接口：`GET /articles/:id/adjacent`。
- 内容服务新增 `AdjacentArticles` 模型，返回上一篇/下一篇文章详情。
- 仓储层按公开文章排序规则查找当前文章前后项，供文章详情页导航使用。

## 容器与部署

- `frontend/nginx/default.conf` 增加 `client_max_body_size 16m`，避免较小图片上传也被 nginx 拦截成 `413 Request Entity Too Large`。
- nginx 增加 `/blog/categories` 与 `/blog/tags` 的静态路由映射。
- `deploy/docker/compose.yaml` 暴露本地 MySQL 端口 `3307:3306`，方便本地调试连接，不改变容器内 MySQL 数据目录。
- `Makefile` 新增/调整部署入口：
  - `micro-redeploy`：本地只重建前端容器。
  - `deploy-static` / `deploy-frontend`：上传静态资源并重建前端。
  - `deploy-code`：服务器拉取指定分支后只重建应用服务。
- README 补充了本地重建、服务器更新、只重建应用容器、不碰 MySQL 数据卷的说明。
- `.gitignore` 忽略 `backups/`、`personal-web-static-update.tar.gz`、`personal-web-release-*.tar.gz`，避免误提交数据库备份或本地发布包。

## 本轮打包上线方式

由于用户希望不通过服务器 `git pull`，而是从本地上传资源到服务器，本轮准备了本地发布包：

```powershell
git archive --format=tar.gz -o personal-web-release-97463f2.tar.gz codex/blog-upload-pet-deploy-20260720
scp .\personal-web-release-97463f2.tar.gz root@117.72.95.156:/tmp/personal-web-release-97463f2.tar.gz
```

服务器侧建议只解压代码并重建应用容器：

```bash
cd /opt/personal-web
tar -xzf /tmp/personal-web-release-97463f2.tar.gz -C /opt/personal-web
docker compose -f deploy/docker/compose.yaml up -d --build frontend media-service content-service web-bff gateway
```

也可以在本地 Docker Desktop 连接远程 Docker context 后执行 compose，但要确认 context 指向服务器，且当前目录是本地最新仓库。更稳妥的方式仍然是 `scp + ssh`。

## 数据安全边界

- 本轮不需要导入 SQL。
- 本轮不需要执行数据库迁移。
- 本轮不需要删除、重建或清空 MySQL volume。
- 不要执行 `docker compose down -v`。
- 如果为了修复历史上传图片访问 403，只允许规范上传文件目录权限，不应改动 MySQL 数据。

## 本地验证结果

已经执行并通过：

```bash
go test . -count=1
go test ./services/media-service/biz/upload ./services/media-service/internal/service -count=1
node --check static/admin/js/admin.js
git diff --cached --check
```

## 上线后建议验证

```bash
curl -I --max-time 10 http://127.0.0.1:8080/blog/
curl -s --max-time 10 http://127.0.0.1:8080/blog/ | grep -E 'blog.css|pet.js'
curl -s --max-time 10 http://127.0.0.1:8080/admin/ | grep 'admin.js?v=13'
```

如果要验证封面上传，建议在后台选择一张图片上传，观察网络请求是否按顺序出现：`init`、若干 `chunks/{index}`、`complete`，并确认 `complete` 返回的 `sha256` 不为空且封面 URL 可访问。

## 未完成/后续可选项

- 当前上传封面时 `biz_id` 可以为空，因为文章可能还没有创建完成；本轮没有实现“文章保存后反向更新媒体记录 biz_id”。
- 素材库图片选择已接入前端入口，但素材来源稳定性和跨域表现仍建议后续继续观察。
- 宠物动作闪烁问题已做运行时层面的缓解，但不同资源帧表质量仍可能影响最终表现，后续可继续按资源逐只校准。
