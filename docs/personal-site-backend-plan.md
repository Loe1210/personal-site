# 个人小站后端开发规划

## 1. 目标

这份规划基于你当前工作区里的 `tiktok_demo` 结构来设计，目标不是重新发明一套架构，而是把它收敛成更适合个人小站的单体后端。

核心原则：

1. 先做单体，不急着拆微服务
2. 目录边界要清晰，后续方便迁移到 `Kitex`
3. 第一版先跑通内容闭环：登录、写文章、发文章、看文章
4. 第二版再考虑缓存、对象存储、AI 能力

## 2. 从 tiktok_demo 可以直接借鉴什么

你这个 `tiktok_demo` 已经给了一个很好的 Hertz 项目骨架：

- `main.go` 负责初始化和启动服务
- `biz/dal` 负责数据层初始化
- `biz/router` 负责路由注册
- `biz/handler` 负责请求绑定和响应
- `biz/service` 负责业务逻辑
- `biz/mw` 负责中间件或基础组件
- `pkg` 负责配置、错误码、工具函数、常量

从你刚才看的代码链路里，可以总结成这条主线：

```text
main -> Init -> dal.Init -> router -> handler -> service -> db
```

这条主线非常适合直接复用到个人小站里。

## 3. 后端总体定位

这个后端的职责有三类：

1. 页面支撑
   - 首页
   - Blog 列表页
   - Blog 详情页
   - About 页

2. 内容管理
   - 管理员登录
   - 文章增删改查
   - 草稿/发布
   - 标签和分类管理
   - 图片上传

3. 工程能力
   - 配置加载
   - 鉴权
   - 统一响应
   - 错误码
   - 日志
   - 分页和查询封装

## 4. 推荐项目结构

建议你新项目最终收敛成这样：

```text
personal_site/
├── biz/
│   ├── handler/
│   │   ├── home/
│   │   ├── blog/
│   │   ├── about/
│   │   ├── admin/
│   │   ├── auth/
│   │   └── upload/
│   ├── service/
│   │   ├── home/
│   │   ├── article/
│   │   ├── tag/
│   │   ├── category/
│   │   ├── auth/
│   │   └── upload/
│   ├── dal/
│   │   ├── db/
│   │   └── init.go
│   ├── model/
│   │   ├── entity/
│   │   ├── dto/
│   │   └── view/
│   ├── router/
│   │   ├── site/
│   │   ├── admin/
│   │   └── register.go
│   └── mw/
│       ├── jwt/
│       ├── logger/
│       └── recover/
├── pkg/
│   ├── configs/
│   ├── constants/
│   ├── errno/
│   ├── utils/
│   └── response/
├── templates/
├── static/
├── migrations/
├── main.go
└── go.mod
```

这个结构和 `tiktok_demo` 的差别不大，所以你学习迁移成本会很低。

## 5. 模块划分

### 5.1 首页模块 home

职责：

- 首页渲染
- 首页基础信息聚合
- 返回站点介绍、入口信息、可选最近文章

对应目录：

- `biz/handler/home`
- `biz/service/home`

### 5.2 内容模块 article

职责：

- 文章列表
- 文章详情
- 后台文章创建
- 后台文章编辑
- 草稿发布状态管理
- slug 生成与唯一性校验

这是你的核心模块，第一阶段要优先完成。

### 5.3 标签模块 tag

职责：

- 标签创建
- 标签列表
- 文章绑定标签
- 标签筛选文章

第一版可以先做简单关系表。

### 5.4 分类模块 category

职责：

- 分类创建
- 分类列表
- 文章关联分类

如果你想简化，分类可以晚于标签做。

### 5.5 认证模块 auth

职责：

- 管理员登录
- JWT 签发与校验
- 管理后台权限保护

这个模块可以直接借鉴 `tiktok_demo` 里的 `jwt` 使用思路，但要简化成单管理员场景。

### 5.6 上传模块 upload

职责：

- 图片上传
- 返回可访问 URL
- 记录文件元数据

第一版建议本地存储，后续再切 `MinIO`。

### 5.7 About 模块 about

职责：

- About 页面内容输出
- 可以先静态配置
- 后续可转成后台可编辑配置项

## 6. 数据库设计建议

第一版建议这些表：

### 6.1 users

```text
id
username
password_hash
nickname
created_at
updated_at
```

说明：

- 第一版只需要一个管理员账号
- 不做复杂用户系统

### 6.2 articles

```text
id
title
slug
summary
content_md
content_html
status
cover_image
category_id
created_at
updated_at
published_at
```

说明：

- `status` 建议：`draft` / `published`
- `content_md` 保存原始 Markdown
- `content_html` 可选缓存渲染结果

### 6.3 tags

```text
id
name
slug
created_at
updated_at
```

### 6.4 article_tags

```text
article_id
tag_id
```

### 6.5 categories

```text
id
name
slug
created_at
updated_at
```

### 6.6 uploads

```text
id
file_name
file_path
file_url
mime_type
size
created_at
```

如果你想再精简，第一版可以先不做 `categories`，只做 `tags`。

## 7. 请求链路怎么设计

参考 `tiktok_demo` 里的 handler 写法，你的请求处理建议固定成下面这个模式：

### 7.1 handler 层职责

- 绑定参数
- 做基础校验
- 调 service
- 统一返回 JSON 或页面渲染

不要把业务逻辑塞进 handler。

### 7.2 service 层职责

- 编排业务逻辑
- 调用 db 方法
- 处理 slug、发布时间、状态切换、标签关联
- 处理 Markdown 转 HTML

### 7.3 db 层职责

- 只做数据库访问
- 不写页面和业务逻辑
- 一类实体一个文件或一组文件

对应你可以直接模仿 `tiktok_demo` 的方式：

```text
handler: 处理请求
service: 处理业务
db: 查询和写库
```

## 8. 路由规划

### 8.1 站点页面路由

```text
GET  /
GET  /blog
GET  /blog/:slug
GET  /about
```

### 8.2 公开 API

```text
GET  /api/articles
GET  /api/articles/:slug
GET  /api/tags
GET  /api/categories
```

### 8.3 管理后台 API

```text
POST   /api/admin/login
GET    /api/admin/articles
POST   /api/admin/articles
PUT    /api/admin/articles/:id
DELETE /api/admin/articles/:id
GET    /api/admin/tags
POST   /api/admin/tags
GET    /api/admin/categories
POST   /api/admin/categories
POST   /api/admin/upload
```

### 8.4 路由分组建议

```text
/site
/api
/api/admin
```

如果你要完全贴近 `tiktok_demo` 的风格，可以每个模块单独建 `router/blog`、`router/about`、`router/admin`。

## 9. 初始化顺序

参考 `tiktok_demo/main.go`，建议你的启动流程写成：

```text
InitConfig()
InitDB()
InitJWT()
InitUploadStore()
InitTemplateRenderer()
RegisterRouter()
RunServer()
```

更接近代码层面的结构：

```text
main.go
  -> Init()
      -> dal.Init()
      -> jwt.Init()
      -> upload.Init()
  -> register(h)
  -> h.Spin()
```

## 10. 中间件规划

### 10.1 JWT 中间件

用途：

- 保护后台接口
- 校验登录状态

### 10.2 日志中间件

用途：

- 记录请求路径、耗时、状态码

### 10.3 恢复中间件

用途：

- 避免 panic 直接把服务打挂

### 10.4 可选：限流中间件

第一版可以不做。

## 11. 配置规划

建议配置文件至少包含：

```text
server:
  port:

database:
  dsn:

jwt:
  secret:
  expire:

upload:
  dir:
  base_url:

site:
  name:
  github_url:
```

放在：

```text
pkg/configs/
```

## 12. 错误码与统一响应

这个部分可以直接学 `tiktok_demo/pkg/errno` 和 `pkg/utils/resp` 的思路。

建议你保留：

- `pkg/errno`
- `pkg/response` 或 `pkg/utils/resp`

统一响应示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

页面渲染接口则单独返回模板。

## 13. 第一阶段开发顺序

### 阶段 1：跑通项目骨架

目标：项目能启动、能访问首页、能连数据库

任务：

1. 初始化 Hertz 项目
2. 搭目录结构
3. 接入 MySQL 和 GORM
4. 配置加载
5. 基础路由注册
6. 首页和 About 静态页返回

### 阶段 2：完成文章主链路

目标：能发文章、能看文章

任务：

1. 建 `articles` 表
2. 做文章列表接口
3. 做文章详情接口
4. 做后台新增文章接口
5. 做后台编辑文章接口
6. 做草稿/发布状态

这是最关键的一阶段。

### 阶段 3：完成后台登录

目标：后台接口受保护

任务：

1. 建 `users` 表
2. 初始化管理员账号
3. 登录接口
4. JWT 中间件
5. 后台接口鉴权

### 阶段 4：补全内容管理能力

目标：后台更好用

任务：

1. 标签管理
2. 分类管理
3. 图片上传
4. 文章筛选和分页
5. 统一错误码与日志

### 阶段 5：工程化完善

目标：项目更稳定、后续更好演进

任务：

1. Docker 化
2. 环境配置拆分
3. 测试补齐
4. 上传能力替换为 `MinIO`
5. Redis 缓存可选接入

## 14. 和 Kitex 的衔接方式

虽然现在不拆微服务，但你从第一天起就可以按领域写清边界。

后面最容易拆出去的模块：

1. `article`
2. `auth`
3. `upload`

也就是说，你现在的 service 边界要像这样：

- `ArticleService`
- `AuthService`
- `UploadService`

以后从“进程内调用”变成 “RPC 调用”时，改动会小很多。

## 15. 和 Eino 的衔接方式

`Eino` 不应该进入第一版主链路。

更适合第二阶段之后接入的功能：

1. 自动生成文章摘要
2. 自动推荐标签
3. 根据文章生成学习卡片
4. 做站内问答

也就是说，`Eino` 更像增强层，不是基础设施层。

## 16. 你现在最该优先写的几个文件

如果下一步正式开工，我建议优先从这些文件入手：

```text
main.go
biz/dal/init.go
biz/dal/db/init.go
biz/router/register.go
biz/handler/home/home_handler.go
biz/handler/blog/blog_handler.go
biz/handler/auth/auth_handler.go
biz/service/article/article_service.go
biz/service/auth/auth_service.go
pkg/configs/
pkg/errno/
```

## 17. 最终建议

一句话总结：

**你的个人小站后端，最适合按 `tiktok_demo` 的分层方式，做成一个边界清晰的 Hertz 单体项目，先完成内容管理闭环，再为 `Kitex` 和 `Eino` 留出演进空间。**

现在不要追求拆服务，先把：

- 登录
- 发文章
- 看文章
- 上传图片
- 标签管理

这些核心链路做扎实。
