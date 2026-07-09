# Phase 02

## 目标

完成单体项目的基础工程骨架，并打通第一个真实业务闭环。

这一阶段重点不在数据库，也不在完整内容系统，而在于先让项目具备继续开发的基础能力。

## 完成内容

1. 初始化 Go 项目并确定 module 为 `github.com/Loe1210/personal-site`
2. 建立了 `biz/`、`pkg/`、`templates/`、`static/` 等基础目录
3. 完成了 `article.proto`、`auth.proto`、`upload.proto` 的 Go 代码生成
4. 建立了最小可运行入口，包括 `main.go`、基础路由注册和首页/健康检查接口
5. 增加了统一响应结构 `pkg/response`
6. 增加了基础错误码 `pkg/errno`
7. 增加了 JWT 工具能力，包括 token 生成和解析
8. 实现了 `POST /api/admin/login`
9. 实现了 `GET /api/admin/me`
10. 接入了 Swagger 文档生成
11. 挂载了 Swagger UI 页面，可直接在浏览器里调试接口

## 设计决策

### 1. proto 只负责模型和接口契约

当前阶段没有把整个项目骨架都交给生成器，而是：

- 用 `proto` 生成模型层
- 手写 `router / handler / service / middleware`

这样更适合理解 Hertz 单体项目的最小链路，也更利于后续按自己的结构演进。

### 2. 登录逻辑不使用库内置 handler

项目中的登录逻辑改为手写 handler，再结合自定义 JWT 中间件处理认证。这样认证模块能和其他模块保持一致风格，也更方便后续服务化。

### 3. Swagger 作为主要调试入口

当前阶段不依赖 Postman，而是把 Swagger 作为主要接口调试工具。这样接口开发和文档沉淀可以同步进行。

## 遇到的问题

### 1. proto 生成路径不符合预期

一开始生成结果落到了错误目录，后来通过调整 `go_package` 和生成参数，最终统一到了 `biz/model/...` 结构下。

### 2. Swagger 文档可生成但静态访问路径不对

一开始 `/docs/swagger.json` 访问失败，后续改为显式暴露 `swagger.json` 和 `swagger.yaml` 文件，并手写一个 Swagger UI 页面，最终调通。

### 3. Swagger Bearer 授权方式容易填错

在 `/api/admin/me` 测试时，必须使用：

`Bearer <token>`

而不是只填 token 本身。

## 当前结果

当前项目已经具备以下基础能力：

- 可启动的 Hertz 服务
- 统一响应结构
- JWT 登录闭环
- 浏览器可调试的 Swagger 文档
- 后续继续开发文章、上传、标签模块的基本地基

## 下一步

下一阶段优先进入文章模块，建议顺序：

1. 设计并实现文章创建接口
2. 设计并实现文章列表接口
3. 设计并实现文章详情接口
4. 将文章接口继续接入 Swagger
5. 在完成文章最小闭环后再接数据库
