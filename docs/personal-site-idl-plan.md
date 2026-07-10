# 个人小站 IDL 驱动开发方案

## 1. 目标

这份方案用于明确：

1. 这个个人小站如何在单体阶段引入 `IDL/thrift`
2. 哪些模块应该先拆成多个 `thrift`
3. 后续开发时每一层需要准备什么
4. 如何避免多个 `thrift` 在生成代码时互相覆盖

这套方案的核心思路是：

**现在仍然做 Hertz 单体，但接口契约按未来微服务边界提前拆好。**

## 2. 总体原则

### 2.1 为什么现在就上 IDL

你已经明确希望：

- 单体阶段就统一接口模型
- 后续学习 `Kitex` 时可以直接复用 `thrift`
- 模块边界从第一天就清晰

所以这个项目不再走“先手写 dto，后补 proto”的路线，而是走：

```text
先设计领域接口
-> 写 thrift
-> 生成模型/handler 基础代码
-> 自己实现 handler/service
-> 单体跑通
-> 后续迁移 Kitex
```

### 2.2 什么坚持自己写

虽然使用 IDL，但下面这些仍然由你自己掌控：

- 路由路径
- handler 实现
- Session 登录逻辑
- service 业务逻辑
- dal/db 访问
- 中间件挂载

也就是说：

**IDL 管契约，不取代你对业务链路的控制。**

## 3. 推荐 thrift 拆分

第一版建议拆成 4 到 6 个 thrift。

### 3.1 `article.thrift`

职责：

- 文章创建
- 文章更新
- 文章删除
- 文章详情
- 文章列表
- 文章发布状态切换

这是第一优先级最高的 thrift。

### 3.2 `auth.thrift`

职责：

- 登录
- 获取当前登录用户
- 可选：刷新 token

因为你已经决定登录逻辑自己写 handler，并且当前登录态已切到 Session，所以 `auth.thrift` 现在用于约束登录、当前用户和退出登录接口。

### 3.3 `upload.thrift`

职责：

- 上传结果描述
- 文件元信息返回
- 可选：上传记录创建

这个模块后续很容易拆成单独服务，所以值得现在就独立。

### 3.4 `tag.thrift`

职责：

- 标签列表
- 标签创建
- 文章标签绑定

### 3.5 `category.thrift`

职责：

- 分类列表
- 分类创建

如果第一版不做分类，这个 thrift 可以晚一点补。

## 4. 哪些内容不建议 thrift 化

下面这些在第一版可以不单独做 proto：

### 4.1 首页渲染

原因：

- 首页更偏页面展示
- 结构变化会比较快
- 不属于未来稳定的服务边界

### 4.2 About 页渲染

原因：

- 偏展示内容
- 更像静态内容或简单配置

### 4.3 模板层 view model

原因：

- 页面模板数据不等于领域契约
- 不适合为了模板渲染把所有 view 层都塞进 proto

所以建议：

- 领域接口使用 proto
- 模板渲染数据可以在 handler/service 层自行组装

## 5. IDL 目录设计

建议目录：

```text
idl/
  article.proto
  auth.proto
  upload.proto
  tag.proto
  category.proto
```

如果后面想统一公共消息，也可以加：

```text
idl/common.proto
```

但第一版不强求，避免过早抽象。

## 6. Go 生成目录设计

为了避免相互覆盖，建议每个 proto 对应一个独立生成目录：

```text
biz/model/
  article/
  auth/
  upload/
  tag/
  category/
```

如果你后续还会生成 http 相关代码，可以保持模块对齐：

```text
biz/handler/
  article/
  auth/
  upload/
  tag/
  category/
```

```text
biz/router/
  article/
  auth/
  upload/
  tag/
  category/
```

## 7. 如何避免多个 proto 互相覆盖

这部分是关键。

### 7.1 每个 proto 都要有独立 `package`

例如：

```proto
package article;
```

```proto
package auth;
```

### 7.2 每个 proto 都要有独立 `go_package`

例如：

```proto
option go_package = "personal_site/biz/model/article";
```

```proto
option go_package = "personal_site/biz/model/auth";
```

这样生成文件天然分目录，不会混在一起。

### 7.3 不同 proto 不要生成到同一个 Go 包

不要这样：

```proto
option go_package = "personal_site/biz/model";
```

如果多个 proto 都写到这个目录里，就很容易冲突。

### 7.4 公共 message 不要重复定义

例如 `BaseResp`、`PageInfo`、`UserSummary` 这类通用结构：

- 要么每个 proto 各自独立，命名足够明确
- 要么后续稳定后再抽到 `common.proto`

第一版更推荐：

- 先局部定义
- 等字段稳定再抽公共

### 7.5 不同 proto 中避免同名 service + 同名文件落到同一路径

例如：

- `ArticleService`
- `AuthService`
- `UploadService`

保持 service 名称和 proto 文件名一致，最不容易混乱。

## 8. 推荐 proto 命名策略

建议使用“领域名 + 动作”的形式，而不是页面名。

### 8.1 `article.proto`

适合的 rpc：

- `CreateArticle`
- `UpdateArticle`
- `DeleteArticle`
- `GetArticleByID`
- `GetArticleBySlug`
- `ListArticles`
- `PublishArticle`

### 8.2 `auth.proto`

适合的 rpc：

- `UserLogin`
- `GetCurrentUser`
- `RefreshUserToken`

### 8.3 `upload.proto`

适合的 rpc：

- `CreateUploadRecord`
- `CompleteUpload`
- `GetUploadInfo`

### 8.4 `tag.proto`

适合的 rpc：

- `CreateTag`
- `ListTags`
- `BindArticleTags`

### 8.5 `category.proto`

适合的 rpc：

- `CreateCategory`
- `ListCategories`

## 9. 单体阶段的请求链路

即使使用 proto，单体里的实际链路仍然建议固定成：

```text
route
-> handler
-> service
-> dal/db
```

只是现在：

- handler 的入参/出参结构体尽量来自 proto 生成结果
- service 的边界也按 proto 对应领域组织

例如：

```text
/api/admin/login
  -> auth handler
  -> auth service
  -> user db
  -> jwt sign
```

```text
/api/admin/articles
  -> article handler
  -> article service
  -> article db
```

## 10. JWT 在这套方案里的位置

你已经决定：

- 不用库内置 `LoginHandler`
- 登录逻辑自己写

所以 JWT 在项目里的角色应该是：

### 10.1 handler 负责

- 接收登录请求
- 调用认证服务
- 登录成功后生成 token
- 返回 token

### 10.2 middleware 负责

- 从请求头解析 token
- 校验 token
- 把用户信息放进上下文

### 10.3 proto 负责

- 定义 `LoginRequest`
- 定义 `LoginResponse`
- 定义 `GetCurrentUserResponse`

所以：

**JWT 和 auth.proto 不冲突。**

## 11. 我后续开发需要用到的东西

这部分是给后面正式开工用的。

### 11.1 IDL 文件

最少需要：

- `idl/article.proto`
- `idl/auth.proto`
- `idl/upload.proto`

可选追加：

- `idl/tag.proto`
- `idl/category.proto`

### 11.2 生成后的模型目录

需要：

```text
biz/model/article/
biz/model/auth/
biz/model/upload/
biz/model/tag/
biz/model/category/
```

### 11.3 路由模块

需要：

```text
biz/router/register.go
biz/router/article/
biz/router/auth/
biz/router/upload/
biz/router/tag/
biz/router/category/
biz/router/site/
```

### 11.4 handler 模块

需要：

```text
biz/handler/article/
biz/handler/auth/
biz/handler/upload/
biz/handler/tag/
biz/handler/category/
biz/handler/site/
```

说明：

- 登录、当前用户、文章 CRUD 都自己写 handler
- 不依赖库内置登录 handler

### 11.5 service 模块

需要：

```text
biz/service/article/
biz/service/auth/
biz/service/upload/
biz/service/tag/
biz/service/category/
```

### 11.6 数据层

需要：

```text
biz/dal/init.go
biz/dal/db/init.go
biz/dal/db/article.go
biz/dal/db/user.go
biz/dal/db/tag.go
biz/dal/db/category.go
biz/dal/db/upload.go
```

### 11.7 中间件

需要：

```text
biz/mw/jwt/
biz/mw/logger/
biz/mw/recover/
```

其中 `jwt` 至少要有：

- token 生成
- token 解析
- 登录态中间件

### 11.8 配置层

需要：

```text
pkg/configs/
```

至少包含：

- server port
- mysql dsn
- jwt secret
- jwt expire
- upload dir
- site basic info

### 11.9 错误码和响应封装

需要：

```text
pkg/errno/
pkg/response/
pkg/utils/
```

至少要解决：

- 统一错误码
- 统一 JSON 响应
- 密码加密
- slug 生成
- 时间格式化

### 11.10 数据表

后续开发至少会用到：

- `users`
- `articles`
- `tags`
- `article_tags`
- `categories`
- `uploads`

### 11.11 页面渲染相关

虽然这不是 proto 的重点，但后续开发仍然需要：

```text
templates/
static/
```

因为首页、Blog 页面和 About 页面都要落地。

## 12. 第一阶段的正式开发顺序

我建议后续按这个顺序开发。

### 第一步：确定 IDL

先完成：

- `article.proto`
- `auth.proto`
- `upload.proto`

如果你想控制范围，先不写 `tag.proto` 和 `category.proto` 也可以。

### 第二步：确定生成目录和命名

先把：

- `package`
- `go_package`
- 生成目录

全部定死，避免后面返工。

### 第三步：搭 Hertz 单体骨架

包括：

- `main.go`
- `dal/init`
- `router/register`
- `jwt middleware`
- 基础配置

### 第四步：先做认证和文章主链路

优先做：

- 用户登录
- 当前用户获取
- 文章创建
- 文章编辑
- 文章列表
- 文章详情

### 第五步：补上传和标签

然后补：

- 图片上传
- 标签管理
- 分类管理

## 13. 最终建议

一句话总结：

**这个项目完全可以在单体阶段使用多个 proto，并按未来微服务边界组织代码；只要 `package`、`go_package` 和生成目录规划清楚，就不会互相覆盖。**

而我后续正式开发最需要的，就是：

1. 明确的 `idl/` 目录
2. 每个 proto 的生成目标目录
3. 自己写的 handler 与 JWT 中间件
4. 按 proto 边界组织的 service 和 dal

这四件事一旦定下来，后面的开发就会非常顺。

## 14. 接口定义建议

这一节把接口写明确，方便后续直接进入 `proto` 设计和代码开发。

说明：

- `rpc` 名称用于 `proto`
- `HTTP` 路径用于 Hertz 路由
- 当前阶段仍然是单体，所以 `HTTP` 和 `rpc` 是“一一对应的领域接口”，不是说一定自动生成路由

### 14.1 `article.proto`

建议 service：

```proto
service ArticleService {}
```

建议接口：

#### `CreateArticle`

- `HTTP`：`POST /api/admin/articles`
- 用途：后台新建文章

请求关键字段：

- `title`
- `slug`
- `summary`
- `content_md`
- `cover_image`
- `category_id`
- `tag_ids`
- `status`

响应关键字段：

- `article`
- `message`

#### `UpdateArticle`

- `HTTP`：`PUT /api/admin/articles/:id`
- 用途：后台编辑文章

请求关键字段：

- `id`
- `title`
- `slug`
- `summary`
- `content_md`
- `cover_image`
- `category_id`
- `tag_ids`
- `status`

响应关键字段：

- `article`
- `message`

#### `DeleteArticle`

- `HTTP`：`DELETE /api/admin/articles/:id`
- 用途：后台删除文章

请求关键字段：

- `id`

响应关键字段：

- `success`
- `message`

#### `GetArticleByID`

- `HTTP`：`GET /api/admin/articles/:id`
- 用途：后台获取文章详情，用于编辑页回填

请求关键字段：

- `id`

响应关键字段：

- `article`

#### `GetArticleBySlug`

- `HTTP`：`GET /api/articles/:slug`
- 页面路由可映射到：`GET /blog/:slug`
- 用途：前台文章详情

请求关键字段：

- `slug`

响应关键字段：

- `article`

#### `ListArticles`

- `HTTP`：`GET /api/articles`
- 页面路由可映射到：`GET /blog`
- 用途：前台文章列表

请求关键字段：

- `page`
- `page_size`
- `tag`
- `category`
- `keyword`
- `status`

说明：

- 前台默认只查 `published`
- 后台列表可以走同一个 rpc，也可以拆管理版接口

响应关键字段：

- `list`
- `total`
- `page`
- `page_size`

#### `ListAdminArticles`

- `HTTP`：`GET /api/admin/articles`
- 用途：后台文章列表

请求关键字段：

- `page`
- `page_size`
- `status`
- `keyword`

响应关键字段：

- `list`
- `total`
- `page`
- `page_size`

#### `PublishArticle`

- `HTTP`：`PATCH /api/admin/articles/:id/publish`
- 用途：单独切换发布状态

请求关键字段：

- `id`
- `status`

响应关键字段：

- `article`
- `message`

### 14.2 `auth.proto`

建议 service：

```proto
service AuthService {}
```

建议接口：

#### `UserLogin`

- `HTTP`：`POST /api/admin/login`
- 用途：用户登录

请求关键字段：

- `username`
- `password`

响应关键字段：

- `token`
- `user`
- `expires_at`

说明：

- 这里由你自己写 handler
- 登录成功后在 handler 或 service 中签发 JWT

#### `GetCurrentUser`

- `HTTP`：`GET /api/admin/me`
- 用途：获取当前登录管理员信息

请求关键字段：

- `token` 通过 header 传递，不一定出现在 body

响应关键字段：

- `user`

#### `RefreshUserToken`

- `HTTP`：`POST /api/admin/refresh`
- 用途：可选，第一版可不做

请求关键字段：

- `token`

响应关键字段：

- `token`
- `expires_at`

### 14.3 `upload.proto`

建议 service：

```proto
service UploadService {}
```

建议接口：

#### `UploadImage`

- `HTTP`：`POST /api/admin/upload`
- 用途：后台上传文章封面或正文图片

请求关键字段：

- `file`
- `biz_type`

说明：

- `file` 实际走 multipart/form-data
- proto 更多用于返回结构和领域边界表达

响应关键字段：

- `file_id`
- `file_name`
- `file_url`
- `mime_type`
- `size`

#### `GetUploadInfo`

- `HTTP`：`GET /api/admin/uploads/:id`
- 用途：查询上传文件信息

请求关键字段：

- `id`

响应关键字段：

- `upload`

### 14.4 `tag.proto`

建议 service：

```proto
service TagService {}
```

建议接口：

#### `CreateTag`

- `HTTP`：`POST /api/admin/tags`
- 用途：后台创建标签

请求关键字段：

- `name`
- `slug`

响应关键字段：

- `tag`

#### `ListTags`

- `HTTP`：`GET /api/tags`
- 后台也可复用：`GET /api/admin/tags`
- 用途：标签列表

请求关键字段：

- `keyword`

响应关键字段：

- `list`

#### `BindArticleTags`

- `HTTP`：通常不单独开放，由文章创建/更新时一起处理
- 用途：绑定文章和标签关系

请求关键字段：

- `article_id`
- `tag_ids`

响应关键字段：

- `success`

### 14.5 `category.proto`

建议 service：

```proto
service CategoryService {}
```

建议接口：

#### `CreateCategory`

- `HTTP`：`POST /api/admin/categories`
- 用途：后台创建分类

请求关键字段：

- `name`
- `slug`

响应关键字段：

- `category`

#### `ListCategories`

- `HTTP`：`GET /api/categories`
- 后台也可复用：`GET /api/admin/categories`

请求关键字段：

- `keyword`

响应关键字段：

- `list`

## 15. 建议的数据结构层次

为了后续开发顺畅，建议每个 proto 至少对应三类结构思维：

### 15.1 Request

例如：

- `CreateArticleRequest`
- `UserLoginRequest`
- `UploadImageRequest`

### 15.2 Response

例如：

- `CreateArticleResponse`
- `ListArticlesResponse`
- `GetCurrentUserResponse`

### 15.3 Entity Summary

例如：

- `Article`
- `Tag`
- `Category`
- `User`
- `UploadFile`

这样做的好处是：

- proto 本身语义稳定
- 以后迁移到 `Kitex` 时几乎可以直接复用
- 单体里 handler/service 也不容易写乱

## 16. 第一版我建议优先落地的接口

如果我们下一步正式开始开发，我建议先只实现这批最关键接口：

### 必做

- `UserLogin`
- `GetCurrentUser`
- `CreateArticle`
- `UpdateArticle`
- `GetArticleBySlug`
- `ListArticles`
- `ListAdminArticles`
- `UploadImage`

### 第二批

- `DeleteArticle`
- `PublishArticle`
- `CreateTag`
- `ListTags`
- `CreateCategory`
- `ListCategories`

这样做的好处是：

- 先把后台登录和文章主链路跑通
- 上传先可用
- 标签和分类随后补足



