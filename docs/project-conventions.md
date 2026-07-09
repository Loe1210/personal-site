# 个人小站项目规范

## 1. 当前目标

当前阶段的唯一目标是：

**完成一个基于 Hertz 的单体个人小站，并在开发过程中持续记录开发日志，保证后续易于 review、易于回顾、易于演进。**

这个项目当前不追求一次性做全，而是按阶段推进：

1. 先完成单体闭环
2. 再补完管理能力
3. 最后为 `Kitex` 和 `Eino` 预留演进空间

## 2. 开发原则

### 2.1 项目形态

- 当前是单体项目
- 核心接口采用 `thrift` 做 IDL
- `handler/service/dal` 边界必须清晰
- 页面渲染和内容接口可以共存

### 2.2 不做的事

- 不做过早微服务拆分
- 不做过重后台系统
- 不把页面 view model 全部 IDL 化
- 不把所有功能一次性塞进第一版

### 2.3 必须保持的约束

- 核心接口优先走 IDL
- 登录逻辑自己写 handler，不使用库内置登录 handler
- JWT 只负责签发和校验
- 每完成一个阶段都补开发日志
- 每个功能闭环完成后再提交合并
- 用户负责代码开发与本地测试，日志整理与 Git 流程由 Codex 协助完成

## 3. 目录规范

建议最终目录结构：

```text
personal_site/
├── idl/
│   ├── article.thrift
│   ├── auth.thrift
│   ├── upload.thrift
│   ├── tag.thrift
│   └── category.thrift
├── biz/
│   ├── handler/
│   ├── service/
│   ├── dal/
│   │   └── db/
│   ├── model/
│   ├── router/
│   └── mw/
├── pkg/
│   ├── configs/
│   ├── constants/
│   ├── errno/
│   ├── response/
│   └── utils/
├── templates/
├── static/
├── docs/
│   ├── devlog/
│   ├── personal-site-ui-spec.md
│   ├── personal-site-backend-plan.md
│   ├── personal-site-idl-plan.md
│   └── project-conventions.md
└── main.go
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

这条流程是后续默认工作方式。

### 4.2 分支策略

不按“模块”分支，而按“阶段 / 功能闭环”分支。

推荐分支：

- `main`
- `feat/project-bootstrap`
- `feat/idl-article-auth-upload`
- `feat/auth-login-jwt`
- `feat/article-crud`
- `feat/upload-image`
- `feat/tag-category`
- `feat/blog-pages`
- `feat/release-prep`

### 4.3 为什么不按模块分支

因为一个真实功能通常会同时修改：

- `idl`
- `handler`
- `service`
- `dal`
- `router`
- 页面
- 文档

如果按模块分支，后面合并会很乱。

### 4.4 合并原则

- 一个分支只承载一个明确阶段目标
- 阶段未完成前，不急着合并
- 阶段完成时必须同步更新开发日志
- 合并前至少保证本地自测通过
- 合并后主分支保持可继续开发状态

## 5. 提交规范

推荐 commit 格式：

```text
feat: add article thrift
feat: implement admin login with jwt
feat: add article admin crud in memory
refactor: simplify article service flow
fix: align thrift field names in article handlers
docs: update phase 02 devlog
```

## 6. 开发日志规范

开发日志统一放在：

```text
docs/devlog/
```

文件命名：

```text
phase-01.md
phase-02.md
phase-03.md
```

每篇开发日志建议固定结构：

```md
# Phase X

## 目标

## 完成内容

## 设计决策

## 遇到的问题

## 当前结果

## 下一步
```

## 7. IDL 规范

### 7.1 每个领域一个 thrift

第一批 thrift：

- `article.thrift`
- `auth.thrift`
- `upload.thrift`

后续 thrift：

- `tag.thrift`
- `category.thrift`

### 7.2 每个 thrift 独立 package

例如：

```thrift
namespace go article
```

### 7.3 不同 thrift 不生成到同一 Go 包

这是为了避免互相覆盖和命名冲突。

## 8. 第一阶段范围

第一阶段只做三件事：

1. 确定目录和规范
2. 起草第一版 IDL
3. 建立第一篇开发日志

不要在这一阶段直接扩展业务逻辑。

## 9. 下一步顺序

Phase 01 结束后，下一步建议按这个顺序做：

1. 初始化项目骨架
2. 实现 JWT 与认证接口
3. 实现文章公开读接口
4. 实现文章后台管理接口
5. 在完成文章最小闭环后再接数据库

## 10. 当前阶段的成功标准

当前阶段完成后，应当具备：

- 目录规范明确
- 开发日志规范明确
- 第一批 thrift 初稿已存在
- 后续开发可以按文档和 thrift 直接进入实现
