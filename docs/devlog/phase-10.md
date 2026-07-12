# Phase 10 - 前后端全面优化

## Goal

从接手项目开始，统览前后端代码，修复性能隐患、安全问题、接口设计缺陷和用户体验问题。

---

## 零、前端架构迁移 — 从服务端渲染到纯静态 SPA（接手后第一件大事）

### 背景

接手前，项目前端基于 Hertz 的服务端模板渲染（phase-08 / phase-09 实现）：
- `templates/` 目录存放所有页面模板
- `biz/site/handler.go` 包含 home / blog / about / admin 等页面渲染 Handler
- `biz/site/route.go` 注册所有页面路由
- 每次页面跳转都由后端渲染完整 HTML 返回

### 迁移内容

**删除原有服务端渲染层：**
- 删除整个 `templates/` 目录（home / blog / about / admin 所有模板）
- 清空 `biz/site/handler.go` 中所有页面渲染 Handler，仅保留 Swagger UI
- 清空 `biz/site/route.go` 中的页面路由

**新建纯静态前端架构：**
- `static/index.html` — 首页（复用 vno.css 视觉风格）
- `static/blog/index.html` — 博客列表页
- `static/blog/post.html` — 文章详情页
- `static/admin/index.html` — 管理后台单页应用
- `static/404.html` — 404 页面

**新建独立前端资源：**
- `static/blog/css/blog.css` — 博客样式（卡片、骨架屏、toast、TOC 等）
- `static/blog/js/api.js` — API 封装（基于 fetch，含分类/标签缓存）
- `static/blog/js/list.js` — 列表页逻辑（分页、筛选、搜索、骨架屏）
- `static/blog/js/post.js` — 详情页逻辑（Markdown 渲染、代码高亮、TOC、复制按钮）
- `static/admin/css/admin.css` — 后台样式
- `static/admin/js/admin.js` — 后台逻辑（登录、文章 CRUD、分类/标签管理）

**重写 `main.go` 路由：**
- 用 `serveSPA` 处理器统一处理所有非 API 路由
- 根据路径前缀（`blog/`、`admin/`）智能回退到对应 HTML
- 保留路径穿越防护（`filepath.Rel` 校验）

### 迁移收益

- 前后端彻底解耦，后端只提供 API，前端独立部署
- 页面跳转无需后端渲染，体验更流畅
- 静态资源可独立缓存、CDN 加速
- 前端代码可独立迭代，不影响后端

---

## 一、后端 - 数据库层

### 1. 数据库连接池配置 (`dal/db/init.go`)

- GORM 默认未配置连接池参数，生产环境会有连接泄漏问题
- 添加连接池配置：
  - `SetMaxIdleConns(10)` — 最大空闲连接数
  - `SetMaxOpenConns(100)` — 最大打开连接数
  - `SetConnMaxLifetime(time.Hour)` — 连接最大存活时间

### 2. 数据库模型修正 (`dal/db/article.go`)

- Article 模型缺少 `IsTop` 和 `AuthorID` 字段，导致文章排序和作者关联无法正常工作
- 添加字段：
  ```go
  IsTop    int   `gorm:"type:tinyint(1);default:0;index"`
  AuthorID int64 `gorm:"default:0"`
  ```

### 3. GORM 日志配置 (`dal/db/init.go`)

- 开发环境（`APP_DEBUG=true` 或 `GIN_MODE=debug`）开启 SQL 日志（Info 级别），方便调试慢查询
- 生产环境默认 `Warn` 级别，避免日志泄露

---

## 二、后端 - Service 层重构

### 4. 内存分页 → 数据库分页 (`service/article.go`)

- 原实现：`Find(&records)` 查出全部文章，再在 Go 里切片分页，文章多了会 OOM
- 重写 `List` 函数，使用 GORM 的 `Limit()` + `Offset()` 实现数据库级分页
- 代码：`query.Offset(int(offset)).Limit(int(pageSize)).Find(&records)`

### 5. N+1 查询问题修复 (`service/article.go`)

- 原实现：`toArticleModel` 每篇文章都单独查一次标签，列表页查 10 篇文章 = 11 次 DB 查询
- 实现 `getTagsForArticles()` 批量查询函数：
  - 先查出所有文章 ID 对应的 `article_tags` 关联
  - 再批量查询标签信息
  - 列表页 DB 查询从 11 次降为 3 次

### 6. ContentHTML 字段语义修正 (`service/article.go`)

- 原实现：`ContentHTML: req.ContentMd`，直接把 Markdown 原文存在 HTML 字段里
- 修正：`CreateArticle` 中 `ContentHTML` 设为空字符串（前端用 marked.js 渲染 Markdown）
- `UpdateArticle` 中更新内容时同步清空 `ContentHTML`
- 保持字段语义正确

### 7. 移除冗余代码 (`service/article.go`)

- 移除 article.go 中重复的 category 和 tag 相关函数，改为调用独立 service
- 移除不存在的 `articlemodel.ListResponse`、`articlemodel.GetPrevNextResponse`、`articlemodel.PrevNextArticle` 等类型引用
- 移除未使用的 `GetPrevNext` 函数和 `toPrevNextArticle` 函数（模型中无对应类型，前端 `getAdjacentPosts` 已返回空）
- 移除 `toArticleModel` 和 `toArticleDetailModel` 中不存在的 `Category` 字段引用
- 清理未使用的 `categorymodel` 导入

### 8. Category Service 补全 (`service/category.go`)

- 新增 `GetCategory` 函数，供 article service 查询分类信息
- 新增 `GetCategoryRequest` / `GetCategoryResponse` 模型 (`biz/model/category/custom.go`)
- `ListCategories` 返回 `errno.Internal` 而非原始错误

### 9. Tag Service 错误处理 (`service/tag.go`)

- `ListTags` 返回 `errno.Internal` 而非原始 `err`

---

## 三、后端 - Handler 层统一

### 10. 统一错误处理 (所有 Handler)

- 原实现：各 Handler 直接 `c.JSON(consts.StatusXXX, response.Error(...))`，错误码和消息散落
- 统一为：所有 Handler 使用 `response.WriteError(c, appErr)` / `response.WriteSuccess(c, resp)`
- 所有 Service 层返回 `errno.AppError` 类型错误
- Handler 中统一判断 `err.(*errno.AppError)` 并转发

涉及的 Handler 文件：
- `biz/article/handler.go` — ListArticles, GetArticleBySlug, GetArticleByID, ListAdminArticles, CreateArticle, UpdateArticle, DeleteArticle
- `biz/auth/handler.go` — Login, Me, Logout
- `biz/category/handler.go` — ListCategories, ListAdminCategories, CreateCategory, UpdateCategory, DeleteCategory
- `biz/tag/handler.go` — ListTags, ListAdminTags, CreateTag, UpdateTag, DeleteTag
- `biz/upload/handler.go` — UploadImage, GetUploadInfo

### 11. Auth Handler 重写 (`biz/auth/handler.go`)

- 原实现：直接 `c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrorCode, err.Error()))`，错误消息暴露
- 重写为统一错误处理模式
- 移除未使用的 `consts` 导入

### 12. 新增错误码 (`pkg/errno/errno.go`)

- 新增 `InvalidCredentials` 错误码（Code: 20013, HTTP 401）
- 用于登录失败场景，替代直接返回 `err.Error()`

---

## 四、后端 - 路由与接口

### 13. 新增基于 ID 的文章查询接口

- 原问题：文章查询只用 slug，但前端传入的参数名为 id，字段名与参数不匹配
- 新增 `GetArticleByID` Handler (`biz/article/handler.go`)
- 新增 `GetPublicArticleByID` Service 函数 (`service/article.go`)
- 路由注册：`GET /api/articles/id/:id` (`biz/article/route.go`)
- 前端查看详情传入 id 参数，用 id 来查看详情

### 14. 文章详情页 URL 修复

- 原问题：点击对应文章 URL 不对，进不去
- 修复：前端 `api.js` 的 `getPost` 改用 `/articles/id/:id` 接口
- 修复：前端 `list.js` 中文章卡片链接改为 `/blog/post/:id`

### 15. 静态文件路由简化 (`main.go`)

- 原实现：通配符路由 `/*filepath` 已处理 `/blog/post/xxx`，但后面又单独定义 `h.GET("/blog/post/:id"...)`，路由重复，存在永远执行不到的代码
- 原实现：`/assets/*filepath` 单独处理，静态文件逻辑分散
- 简化为单个 `serveSPA` 处理器：
  - 统一处理所有非 API 路由
  - 根据路径前缀（`blog/`、`admin/`）智能回退到对应 HTML
  - 保留路径穿越防护（`filepath.Rel` 校验）
- 移除冗余的重定向路由和 `serveHTML` / `staticFileHandler` 函数

### 16. 添加缓存头中间件 (`main.go`)

- 新增 `staticCacheMiddleware`，为静态资源添加 `Cache-Control: public, max-age=86400`
- CSS/JS/图片/字体不再每次重新请求

---

## 五、后端 - 安全

### 17. 删除默认测试用户

- 移除 `seedTestUser()` 函数及相关代码
- 删除 `editor1 / 123456` 默认测试用户创建逻辑
- 保留 `admin` 主管理员账户

### 18. 移除登录页密码提示

- 从 `static/admin/index.html` 中删除 `<p class="login-hint">默认账号：admin / 密码：hins123</p>` 标签
- 避免敏感信息暴露在页面源码中

---

## 六、前端 - 代码质量

### 19. 代码高亮语言包不全 (`static/blog/post.html`)

- 原实现：只引入 `common.min.js`，只支持常见语言，Go/Python/Rust 等都不高亮
- 替换为 `highlight.min.js`，支持全语言高亮

### 20. 全局变量污染 (`static/blog/js/api.js`)

- 原实现：所有变量（`API_BASE`、`categoryCache`、`MOCK_POSTS`、`BlogAPI` 等）都在全局作用域
- 使用 IIFE 封装，仅通过 `window.BlogAPI` 和 `window.formatDate` 导出必要接口
- `list.js` 和 `post.js` 已使用 IIFE 封装（无需修改）

---

## 七、前端 - 用户体验

### 21. 骨架屏加载状态 (`static/blog/js/list.js`)

- 原实现：只有文字"加载中..."
- 新增骨架屏卡片（6 个占位卡片 + shimmer 动画），视觉体验更好
- 新增 `.post-card--skeleton` 和 `.skeleton-line` CSS 样式 (`static/blog/css/blog.css`)

### 22. 图片懒加载 (`static/blog/js/post.js`)

- 原实现：文章里的图片没有 `loading="lazy"`，长文章打开慢
- 渲染后为所有 `.post-content img` 添加 `loading="lazy"` 属性

### 23. 错误提示友好化 (`static/blog/js/list.js`, `post.js`)

- 原实现：API 请求失败只在控制台打 log，用户看不到任何提示
- 新增 `showToast` 函数，API 请求失败时显示底部 toast 提示
- 列表页和详情页均已接入
- 新增 `.blog-toast` CSS 样式

### 24. 返回顶部按钮

- 已在 `list.js` 和 `post.js` 中实现 `bindBackToTop`（滚动超过 400px 显示）
- HTML 中已有 `#backToTop` 按钮元素和对应 CSS 样式

### 25. 标签/分类点击筛选 (`static/blog/js/list.js`)

- 原实现：文章卡片上显示了标签，但点击标签不能筛选同标签文章
- 文章卡片的标签从 `<span>` 改为 `<a>` 链接，点击跳转到 `/blog/?tag=xxx`
- 文章卡片的分类同样改为可点击链接，点击跳转到 `/blog/?category=xxx`
- 新增 hover 样式

---

## 八、前端 - UI 修复（前期工作）

### 26. 窗口缩小后右上角目录按钮处理

- 问题：窗口缩小后右上角出现目录按钮，悬停时只有黑框看不见内容
- 处理：根据用户要求删除该按钮（撤销了多次尝试修复的更改）

### 27. 滑动黑线问题

- 问题：滑动时偶尔出现黑线
- 状态：尝试修复但未完全解决，用户反馈"还是有黑条"
- 此问题待后续排查（可能与 CSS `transform` 或 `overflow` 有关）

---

## 涉及的文件清单

### 后端
- `dal/db/init.go` — 连接池配置、GORM 日志、移除 seedTestUser
- `dal/db/article.go` — Article 模型添加 IsTop、AuthorID 字段
- `service/article.go` — 完全重写（数据库分页、N+1 修复、统一错误处理、移除冗余代码）
- `service/category.go` — 新增 GetCategory 函数、统一错误处理
- `service/tag.go` — 统一错误处理
- `service/auth.go` — 使用 InvalidCredentials 错误码
- `biz/article/handler.go` — 统一错误处理、新增 GetArticleByID、移除 GetArticlePrevNext
- `biz/article/route.go` — 新增 `/api/articles/id/:id` 路由
- `biz/auth/handler.go` — 重写为统一错误处理模式
- `biz/category/handler.go` — 统一错误处理
- `biz/tag/handler.go` — 统一错误处理
- `biz/upload/handler.go` — 统一错误处理
- `biz/site/handler.go` — 删除所有页面渲染 Handler，仅保留 Swagger
- `biz/site/route.go` — 删除所有页面路由
- `biz/model/category/custom.go` — 新增 GetCategoryRequest/Response
- `pkg/errno/errno.go` — 新增 InvalidCredentials 错误码
- `main.go` — 静态文件路由简化、缓存头中间件、serveSPA

### 前端（新建）
- `static/index.html` — 首页（vno 风格）
- `static/blog/index.html` — 博客列表页
- `static/blog/post.html` — 文章详情页
- `static/admin/index.html` — 管理后台 SPA
- `static/404.html` — 404 页面
- `static/blog/css/blog.css` — 博客样式（卡片、骨架屏、toast、TOC 等）
- `static/blog/js/api.js` — API 封装（IIFE）
- `static/blog/js/list.js` — 列表页逻辑（骨架屏、toast、标签筛选）
- `static/blog/js/post.js` — 详情页逻辑（懒加载、toast、TOC）
- `static/admin/css/admin.css` — 后台样式
- `static/admin/js/admin.js` — 后台逻辑

### 前端（修改）
- `static/admin/index.html` — 移除登录页密码提示
- `static/blog/post.html` — 代码高亮语言包替换

### 删除
- `templates/` 目录（原服务端渲染模板，全部删除）

### 文档
- `docs/devlog/phase-10.md` — 本开发日志

---

## 验证

- `go build` 编译通过
- 所有 Handler 和 Service 层错误处理统一
- 新增 `/api/articles/id/:id` 接口可用
- 前端纯静态架构正常运行
- 骨架屏、toast、懒加载、标签筛选样式已就绪
- 登录页密码提示已移除
- 默认测试用户已删除
