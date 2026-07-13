# 阶段 A / B 与 Swagger 清理实施计划

> **给执行型 agent 的要求：** 必须使用 `superpowers:subagent-driven-development`（推荐）或 `superpowers:executing-plans` 按任务逐项执行。本计划使用 `- [ ]` 复选框跟踪步骤。

**目标：** 在不进入 RBAC 后续工作的前提下，完成阶段 A 的用户侧修复、彻底移除 Swagger，并补齐第一轮阶段 B 工程化基础。

**架构思路：** 保持当前“静态 SPA + Go API”结构不变。先补一个最小后端能力支撑上下篇导航，同时移除 Swagger 运行面；再补配置拆分、Docker 交付和关键回归测试。

**技术栈：** Go、Hertz、GORM、YAML 配置、静态 HTML/CSS/JS、Docker

## 全局约束

- 本计划不实现 P-06 至 P-12 的 RBAC 管理接口。
- 前端继续只在 `static/` 目录迭代，不回退到服务端渲染。
- 保持当前 session 认证模型不变。
- 新增后端接口继续使用现有 `errno` + `response` 的错误处理模式。
- 黑线修复和后台视觉对齐优先采用最小 CSS 改动。
- Swagger 页面、静态规范文件和运行时暴露都要彻底移除。

---

## 文件结构与职责

- 修改：`service/article.go`
  责任：文章查询逻辑，包括新增相邻文章查询能力。
- 修改：`biz/article/handler.go`
  责任：暴露相邻文章 API，并沿用现有响应封装。
- 修改：`biz/article/route.go`
  责任：注册相邻文章路由。
- 修改：`static/blog/js/api.js`
  责任：改为真实调用相邻文章 API，而不是返回空占位。
- 修改：`static/blog/js/post.js`
  责任：仅在存在相邻文章时渲染上下篇导航。
- 修改：`static/blog/css/blog.css`
  责任：修复黑线问题，并顺手收口文章导航样式。
- 修改：`static/admin/css/admin.css`
  责任：后台 Terminal Gallery 风格对齐。
- 修改：`static/admin/index.html`
  责任：如有必要，仅做少量结构或文案微调。
- 修改：`main.go`
  责任：移除 Swagger 暴露，并在清理中顺手收紧缓存作用范围。
- 修改或删除：`biz/site/handler.go`、`biz/site/route.go`
  责任：删除仅剩的 Swagger 站点 handler 能力。
- 修改：`configs/config.go`
  责任：拆分配置模型，同时保留默认值和环境变量覆盖。
- 修改：`configs/config.yaml`
  责任：使本地示例配置与拆分后的结构保持一致。
- 新建：`Dockerfile`
  责任：应用镜像构建。
- 新建：`docker-compose.yml`
  责任：本地联动应用与 MySQL。
- 新建或修改：`README.md` 或运行文档
  责任：记录配置与 Docker 使用方式。
- 新建：聚焦型 `_test.go` 文件
  责任：覆盖配置加载与上下篇导航关键回归逻辑。

## 任务 1：补齐上下篇文章导航

**文件：**
- 修改：`C:\Users\Administrator\Desktop\personal web\service\article.go`
- 修改：`C:\Users\Administrator\Desktop\personal web\biz\article\handler.go`
- 修改：`C:\Users\Administrator\Desktop\personal web\biz\article\route.go`
- 修改：`C:\Users\Administrator\Desktop\personal web\static\blog\js\api.js`
- 修改：`C:\Users\Administrator\Desktop\personal web\static\blog\js\post.js`
- 测试：`C:\Users\Administrator\Desktop\personal web\service\article_adjacent_test.go`

**接口边界：**
- 输入：现有文章模型、公开文章排序规则、统一响应封装
- 输出：`GetAdjacentPublicArticles(ctx context.Context, id int64) (*articlemodel.GetAdjacentArticlesResponse, error)` 与 `GET /api/articles/id/:id/adjacent`

- [ ] **步骤 1：先写失败测试**

```go
func TestGetAdjacentPublicArticlesReturnsNeighborsInPublicOrder(t *testing.T) {
	t.Skip("implement with isolated test DB setup")
}
```

- [ ] **步骤 2：运行测试，确认先失败**

运行：`go test ./service -run TestGetAdjacentPublicArticlesReturnsNeighborsInPublicOrder -v`  
预期：FAIL，因为测试仍是占位或目标函数尚不存在。

- [ ] **步骤 3：补最小实现**

实现要求：

- 取出所有已发布文章的 ID，排序规则与前台公开列表一致
- 找到当前文章所在位置
- 返回最小相邻文章信息：上一条、下一条
- 增加 handler 和 route 暴露该能力

- [ ] **步骤 4：替换前端空占位逻辑**

把 `BlogAPI.getAdjacentPosts` 从：

```js
getAdjacentPosts: function (slugOrId) {
    return Promise.resolve({ prev: null, next: null });
},
```

改成真实请求：

```js
getAdjacentPosts: function (id) {
    return request('/articles/id/' + encodeURIComponent(id) + '/adjacent').then(function (data) {
        return {
            prev: data.prev ? mapBackendArticle(data.prev) : null,
            next: data.next ? mapBackendArticle(data.next) : null
        };
    });
},
```

- [ ] **步骤 5：运行本任务验证**

运行：`go test ./service ./biz/article -v`  
预期：PASS

- [ ] **步骤 6：提交**

```bash
git add service/article.go biz/article/handler.go biz/article/route.go static/blog/js/api.js static/blog/js/post.js service/article_adjacent_test.go
git commit -m "feat: restore adjacent post navigation"
```

## 任务 2：修复滚动黑线并稳定博客样式

**文件：**
- 修改：`C:\Users\Administrator\Desktop\personal web\static\blog\css\blog.css`
- 修改：`C:\Users\Administrator\Desktop\personal web\static\blog\post.html`
- 测试：人工验证记录

**接口边界：**
- 输入：现有博客详情布局与浮层组件
- 输出：滚动过程中无明显黑线或合成缝隙

- [ ] **步骤 1：定位 CSS 热点区域**

重点检查这些选择器及其相关容器：

```css
body
.blog-shell
.post-layout
.post-toc
.admin-entry-btn
#backToTop
```

- [ ] **步骤 2：做最小 CSS 修复**

优先方向：

```css
html, body {
    overflow-x: hidden;
    background: #0b1220;
}
```

再把全页级 `transform` 或易出缝隙的背景承载方式收紧到局部浮层。

- [ ] **步骤 3：人工验证**

本地启动项目，打开博客详情页，在桌面宽度和较窄宽度下滚动检查。  
预期：页面边缘或图层交界处不再出现间歇性黑线。

- [ ] **步骤 4：提交**

```bash
git add static/blog/css/blog.css static/blog/post.html
git commit -m "fix: stabilize blog scrolling visuals"
```

## 任务 3：后台视觉对齐

**文件：**
- 修改：`C:\Users\Administrator\Desktop\personal web\static\admin\css\admin.css`
- 修改：`C:\Users\Administrator\Desktop\personal web\static\admin\index.html`

**接口边界：**
- 输入：现有后台 SPA 结构
- 输出：更统一的 Terminal Gallery 风格，不改变核心交互

- [ ] **步骤 1：审查主要视觉状态**

检查这些状态：

- 登录弹层
- Tab 区域
- 空状态
- 加载态
- 文章列表项操作区
- 分类 / 标签管理卡片

- [ ] **步骤 2：只做样式收口**

聚焦：

- 字体与字重
- 间距体系
- 徽标与状态标签
- 卡片层次
- 对比度与按钮统一性

- [ ] **步骤 3：人工验证**

打开 `/admin/`，登录后检查 posts、categories、tags 三个视图。  
预期：主要界面处于同一视觉语言下，没有明显突兀的组件。

- [ ] **步骤 4：提交**

```bash
git add static/admin/css/admin.css static/admin/index.html
git commit -m "style: align admin ui with site visual system"
```

## 任务 4：彻底移除 Swagger

**文件：**
- 修改：`C:\Users\Administrator\Desktop\personal web\main.go`
- 修改或删除：`C:\Users\Administrator\Desktop\personal web\biz\site\handler.go`
- 修改或删除：`C:\Users\Administrator\Desktop\personal web\biz\site\route.go`
- 删除或停用：`C:\Users\Administrator\Desktop\personal web\docs\swagger.json`
- 删除或停用：`C:\Users\Administrator\Desktop\personal web\docs\swagger.yaml`
- 删除或停用：`C:\Users\Administrator\Desktop\personal web\docs\docs.go`

**接口边界：**
- 输入：当前站点启动逻辑
- 输出：Swagger 不再出现在运行时暴露面，也不再保留无用文档入口

- [ ] **步骤 1：先写一个失败的路由验证**

加一个轻量测试或 handler 级验证，确保 `/swagger.json` 与 `/swagger.yaml` 不再被服务。

- [ ] **步骤 2：运行验证，确认先失败**

运行：`go test ./... -run Swagger -v`  
预期：FAIL，直到相关入口被彻底移除。

- [ ] **步骤 3：移除运行时暴露**

从 `main.go` 删除：

```go
h.StaticFile("/swagger.json", mustAbs(root, "docs", "swagger.json"))
h.StaticFile("/swagger.yaml", mustAbs(root, "docs", "swagger.yaml"))
```

同时删除 Swagger 专用 site handler 代码。

- [ ] **步骤 4：清理仓库遗留产物**

如果没有任何编译或运行依赖，直接删除相关 swagger 文件；若仍有引用，先删依赖再删文件。

- [ ] **步骤 5：运行验证**

运行：`go test ./... -v`  
预期：PASS

- [ ] **步骤 6：提交**

```bash
git add main.go biz/site/handler.go biz/site/route.go docs/docs.go docs/swagger.json docs/swagger.yaml
git commit -m "chore: remove swagger surface"
```

## 任务 5：拆分配置结构

**文件：**
- 修改：`C:\Users\Administrator\Desktop\personal web\configs\config.go`
- 修改：`C:\Users\Administrator\Desktop\personal web\configs\config.yaml`
- 测试：`C:\Users\Administrator\Desktop\personal web\configs\config_test.go`

**接口边界：**
- 输入：`configs.Load`、`configs.AppConfig`、环境变量覆盖逻辑
- 输出：更清晰的配置分组，同时保持当前默认行为兼容

- [ ] **步骤 1：先写失败测试**

覆盖这些场景：

- 配置文件缺失时走默认值
- YAML 值可覆盖默认值
- 环境变量可覆盖 YAML 值

- [ ] **步骤 2：运行测试，确认先失败**

运行：`go test ./configs -run TestLoad -v`  
预期：FAIL，直到新结构与测试预期完成。

- [ ] **步骤 3：实现配置分组扩展**

在保留 `server`、`mysql`、`session` 的同时，新增：

- `upload`
- `site`

并保持当前环境变量覆盖行为继续可用。

- [ ] **步骤 4：更新示例配置**

同步扩展 `configs/config.yaml`，保证本地启动体验不被破坏。

- [ ] **步骤 5：运行验证**

运行：`go test ./configs -v`  
预期：PASS

- [ ] **步骤 6：提交**

```bash
git add configs/config.go configs/config.yaml configs/config_test.go
git commit -m "refactor: split application config domains"
```

## 任务 6：补齐 Docker 交付

**文件：**
- 新建：`C:\Users\Administrator\Desktop\personal web\Dockerfile`
- 新建：`C:\Users\Administrator\Desktop\personal web\docker-compose.yml`
- 新建或修改：`C:\Users\Administrator\Desktop\personal web\README.md`

**接口边界：**
- 输入：应用二进制、静态资源、配置注入方式
- 输出：可复现的容器构建与本地 compose 启动方式

- [ ] **步骤 1：补应用镜像构建**

新建多阶段 `Dockerfile`，要求：

- 编译 Go 二进制
- 复制 `static/` 与运行时仍需保留的配置资源
- 使用稳定的启动命令

- [ ] **步骤 2：补本地 compose**

新建 `docker-compose.yml`，把应用和 MySQL 接起来，并通过环境变量完成覆盖。

- [ ] **步骤 3：补运行说明**

至少记录这些命令：

```bash
docker build -t personal-site .
docker compose up --build
```

- [ ] **步骤 4：运行验证**

运行：`docker build -t personal-site .`  
预期：exit 0

- [ ] **步骤 5：提交**

```bash
git add Dockerfile docker-compose.yml README.md
git commit -m "build: add docker delivery files"
```

## 任务 7：补关键回归测试并更新项目记录

**文件：**
- 修改：前面新增的测试文件
- 修改：`C:\Users\Administrator\Desktop\personal web\docs\pending-tasks.md`
- 修改：`C:\Users\Administrator\Desktop\personal web\docs\devlog\phase-10.md` 或按需要追加新的开发日志

**接口边界：**
- 输入：前 6 个任务的落地结果
- 输出：更新后的待办状态与新鲜验证证据

- [ ] **步骤 1：更新 pending-tasks 状态**

在 `docs/pending-tasks.md` 中标记本轮已完成项，并保留 RBAC 延后状态不变。

- [ ] **步骤 2：记录开发日志**

补一段简洁日志，说明本轮改动、剩余事项与仍需人工验证的风险点。

- [ ] **步骤 3：运行全量验证**

运行：`go test ./...`  
预期：PASS

运行：`go build ./...`  
预期：PASS

如果本机可用，再运行：

运行：`docker build -t personal-site .`  
预期：PASS

- [ ] **步骤 4：提交**

```bash
git add docs/pending-tasks.md docs/devlog/phase-10.md
git commit -m "docs: record phase a and b progress"
```

## 自检

- 规格覆盖：任务 1 至任务 4 覆盖阶段 A 与 Swagger 清理；任务 5 至任务 7 覆盖阶段 B 与文档更新。
- 占位检查：计划中没有 `TBD`、`TODO` 或“后面再补”这种空描述。
- 命名一致性：相邻文章接口在 service、handler、route 和前端消费层保持同名语义一致。
