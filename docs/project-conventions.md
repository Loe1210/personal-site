# 个人小站项目规范

## 1. 当前目标

当前阶段的目标是：

**完成一个可本地运行、可持续演进的 Hertz + 微服务个人小站，并在开发过程中持续记录开发日志，保证后续易于 review、易于回顾、易于演进。**

这个项目当前的工作方式是：

1. 浏览器入口走 `frontend` + `Nginx`
2. 外部 HTTP 统一进入 `gateway`
3. 业务能力拆到 `auth-service`、`content-service`、`media-service`、`web-bff`
4. 认证使用 `session cookie + Redis`
5. 后续 RPC 通过 `proto` 和 `Kitex` 逐步演进

## 2. 开发原则

### 2.1 项目形态

- 当前是微服务项目，不再回到旧单体入口
- 核心接口采用 `proto` 做 IDL
- 每个服务都要有自己的 `internal/model`
- `handler/service/dal` 边界必须清晰
- `frontend` 只负责静态页面和 `/api` 转发，不承载业务逻辑

### 2.2 不做的事

- 不恢复旧单体入口、旧 `biz/`、旧 `service/` 和旧根 IDL
- 不做过早耦合的服务间 RPC
- 不把页面 view model 全部 proto 化
- 不把所有功能一次性塞进第一版

### 2.3 必须保持的约束

- 核心接口优先走 IDL
- 登录逻辑自己写 handler，不使用库内置登录 handler
- Session 通过 cookie + Redis 共享
- 每完成一个阶段都补开发日志
- 每个功能闭环完成后再提交合并

## 3. 目录规范

建议最终目录结构：

```text
personal_site/
├── frontend/
├── services/
│   ├── gateway/
│   ├── auth-service/
│   ├── content-service/
│   ├── media-service/
│   └── web-bff/
├── idl/
│   ├── auth/
│   ├── content/
│   └── media/
├── docs/
│   ├── devlog/
│   ├── runbooks/
│   └── superpowers/
├── deploy/
└── README.md
```

## 4. Git 工作流

当前仓库目标：

- GitHub 仓库：`Loe1210/personal-site`
- 仓库地址：`https://github.com/Loe1210/personal-site.git`

后续统一使用 Git 推进，不再停留在“本地只改不提交”的状态。

### 4.1 标准阶段流程

每完成一个阶段，固定执行以下流程：

1. 本地完成一个阶段
2. 更新对应 `devlog`
3. 提交 `commit`
4. 推送到功能分支
5. 合并到主分支

### 4.2 分支策略

不按“模块”分支，而按“阶段 / 功能闭环”分支。

### 4.3 合并原则

- 一个分支只承载一个明确阶段目标
- 阶段未完成前，不急着合并
- 阶段完成时必须同步更新开发日志
- 合并前至少保证本地自测通过
- 合并后主分支保持可继续开发状态

## 5. 提交规范

推荐 commit 格式：

```text
feat: add article proto
feat: implement admin login with session
docs: add phase 01 development log
refactor: simplify article service flow
fix: handle empty article slug
```

## 6. 开发日志规范

开发日志统一放在：

```text
docs/devlog/
```

## 7. IDL 规范

- 核心接口使用 `proto`
- 每个 proto 独立 `package`
- 每个 proto 独立 `go_package`
- 避免多个 proto 生成到同一个 Go 包

## 8. 当前阶段成功标准

当前阶段完成后，应当具备：

- 浏览器入口明确
- 微服务入口明确
- 目录规范明确
- 开发日志规范明确
- 后续开发可以按文档和 proto 直接进入实现