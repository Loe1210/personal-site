# Phase 02

## 目标

完成单体项目的基础工程骨架，并打通第一个真实业务闭环。

这一阶段重点不在数据库，也不在完整内容系统，而在于先让项目具备继续开发的基础能力，并完成从 `proto` 到 `thrift` 的路线切换。

## 完成内容

1. 初始化 Go 项目并确定 module 为 `github.com/Loe1210/personal-site`
2. 建立了 `biz/`、`pkg/`、`templates/`、`static/` 等基础目录
3. 完成了 `auth.thrift`、`article.thrift`、`upload.thrift` 的初版起草与模型生成
4. 建立了最小可运行入口，包括 `main.go`、基础路由注册和首页/健康检查接口
5. 增加了统一响应结构 `pkg/response`
6. 增加了基础错误码 `pkg/errno`
7. 增加了 JWT 工具能力，包括 token 生成、解析和认证中间件
8. 实现并验证了 `POST /api/admin/login`
9. 实现并验证了 `GET /api/admin/me`
10. 接入了 Swagger 文档生成与 UI 调试页
11. 实现并验证了公开文章接口：`GET /api/articles`、`GET /api/articles/:slug`
12. 实现并验证了后台文章接口：`GET /api/admin/articles`、`POST /api/admin/articles`
13. 实现并验证了后台文章更新与删除接口：`PUT /api/admin/articles/:id`、`DELETE /api/admin/articles/:id`
14. 当前文章模块先采用内存假数据，便于先打通接口闭环再接数据库

## 设计决策

### 1. 从 proto 切换到 thrift

结合后续 `Hertz + Kitex` 的学习和演进方向，项目核心 IDL 从 `proto` 切换为 `thrift`。这样后续无论继续做单体，还是拆成基于 Kitex 的服务，模型和接口契约都更统一。

### 2. 项目骨架继续手写，IDL 只负责契约和模型

当前阶段没有把整个项目骨架都交给生成器，而是：

- 用 `thrift` 生成模型层
- 手写 `router / handler / service / middleware`

这样更适合理解 Hertz 单体项目的最小链路，也更利于后续按自己的结构演进。

### 3. 登录逻辑不使用库内置 handler

项目中的登录逻辑改为手写 handler，再结合自定义 JWT 中间件处理认证。这样认证模块能和其他模块保持一致风格，也更方便后续服务化。

### 4. Swagger 作为主要调试入口

当前阶段不依赖 Postman，而是把 Swagger 作为主要接口调试工具。这样接口开发、联调和文档沉淀可以同步进行。

### 5. 文章模块先跑通最小闭环，再接数据库

文章模块先实现“公开读取 + 后台 CRUD”的内存版闭环，目的不是做最终形态，而是优先把路由、参数绑定、响应结构、JWT 保护和 Swagger 链路跑顺。

## 遇到的问题

### 1. proto 生成路径和包结构不符合预期

一开始生成结果落到了错误目录，后来通过调整 `go_package` 和生成参数，最终统一到了 `biz/model/...` 结构下。随后项目路线切换到 `thrift`，对应模型也重新整理到了按领域分包的结构。

### 2. thriftgo 与 apache/thrift 运行时版本不兼容

`thriftgo 0.4.5` 生成的 Go 代码使用的是旧版 `apache/thrift` 接口风格，而项目最初拉到的是新版运行时，导致 `Skip` 等方法签名不一致。最终通过将 `github.com/apache/thrift` 降到兼容版本并重新生成模型解决。

### 3. Swagger 文档可生成但静态访问路径不对

一开始 `/docs/swagger.json` 访问失败，后续改为显式暴露 `swagger.json` 和 `swagger.yaml` 文件，并手写一个 Swagger UI 页面，最终调通。

### 4. Swagger Bearer 授权方式容易填错

在受保护接口测试时，必须使用：

`Bearer <token>`

而不是只填 token 本身。

### 5. thrift 生成字段名和手写代码容易不一致

在 `auth` 和 `article` 模块迁移过程中，多次出现字段名与生成模型不一致的问题，例如：

- `ID` 不是 `Id`
- `ContentHTML` 不是 `ContentHtml`
- `TagIds` 不是 `TagIDs`

这也进一步说明后续开发时必须先对照生成的 model 再写 handler 和 service。

## 当前结果

当前项目已经具备以下基础能力：

- 可启动的 Hertz 服务
- 统一响应结构
- JWT 登录闭环
- 浏览器可调试的 Swagger 文档
- 基于 thrift 的 auth / article / upload 模型基础
- 公开文章读取能力
- 后台文章完整 CRUD 能力（内存版）
- 后续继续接数据库、上传、标签分类模块的基本地基

## 下一步

下一阶段优先进入数据库接入，建议顺序：

1. 为文章模块设计表结构与初始化方式
2. 接入数据库配置与连接管理
3. 把文章模块从内存数据切到持久化存储
4. 补文章发布状态、分类、标签等关联能力
5. 在文章模块稳定后再进入上传模块与页面整合

