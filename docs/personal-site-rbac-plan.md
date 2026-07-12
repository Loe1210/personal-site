# 个人小站 RBAC 最小闭环方案

## 1. 目标

这一阶段的目标不是一次性做完整后台，而是先把管理端权限体系跑通，保证后续 `article`、`category`、`tag`、`upload` 等接口都能在统一规则下做鉴权。

当前阶段的最小闭环包括：

1. 用户通过 Session 登录
2. 登录用户与角色建立绑定关系
3. 角色与权限点建立绑定关系
4. 管理端路由可以基于权限点做拦截
5. 无权限时返回统一的 `403` 业务响应

## 2. 为什么现在做 RBAC

当前管理端接口已经开始成型：

- `auth`
- `article`
- `category`
- `tag`
- 后续 `upload`

如果权限模型不先定下来，后面每新增一个后台接口都要重新补权限判断，返工成本会越来越高。

因此这一阶段先把权限骨架搭好，后面新增后台接口时只需要补权限点和路由绑定，不需要再重做鉴权模型。

## 3. 最小领域拆分

当前建议把 RBAC 独立成一个领域，而不是继续塞进 `auth`：

- `auth` 负责登录、session 登录态、当前用户识别
- `rbac` 负责角色、权限、绑定关系、权限校验

这样后续即使把认证和权限拆成单独服务，边界也比较清晰。

## 4. 数据表设计

### 4.1 users

当前登录流程已经切到数据库用户，后续继续围绕 `users` 扩展即可。

核心字段建议：

- `id`
- `username`
- `password_hash`
- `nickname`
- `type`
- `status`
- `created_at`
- `updated_at`

说明：

- `password_hash` 必须存加密后的密码，不存明文
- `status` 用于禁用账号
- `type` 用于区分后台用户等身份类型

### 4.2 roles

角色表：

- `id`
- `name`
- `code`
- `description`
- `created_at`
- `updated_at`

建议约束：

- `code` 全局唯一，例如 `super_admin`、`editor`

### 4.3 permissions

权限点表：

- `id`
- `name`
- `code`
- `resource`
- `action`
- `description`
- `created_at`
- `updated_at`

建议示例：

- `article:create`
- `article:update`
- `article:delete`
- `category:create`
- `tag:create`
- `upload:create`

说明：

- `code` 是真正用于程序鉴权的关键字段
- `resource + action` 是便于理解和扩展的辅助拆分

### 4.4 user_roles

用户和角色的关联表：

- `id`
- `user_id`
- `role_id`
- `created_at`

说明：

- 采用多对多，避免未来只能一个用户一个角色

### 4.5 role_permissions

角色和权限的关联表：

- `id`
- `role_id`
- `permission_id`
- `created_at`

说明：

- 一个角色可拥有多个权限点
- 一个权限点也可复用于多个角色

## 5. 第一版角色建议

第一版保留两个角色：

### 5.1 super_admin

拥有全部后台权限。

用途：

- 你自己的后台账号
- 初期开发和演示账号

### 5.2 editor

拥有内容管理相关权限，但不拥有角色权限配置能力。

用途：

- 用于验证低权限用户的真实 `403` 场景
- 后续也可继续作为编辑角色演进

## 6. 第一版权限点建议

建议先只覆盖已经存在或即将开发的后台接口：

- `article:read`
- `article:create`
- `article:update`
- `article:delete`
- `category:read`
- `category:create`
- `tag:read`
- `tag:create`
- `user:me`
- `user:logout`
- `upload:create`
- 后续再补 RBAC 管理接口权限

## 7. 中间件设计

当前已经拆成两层：

### 7.1 AuthMiddleware

只负责：

1. 读取 session
2. 判断当前用户是否已登录
3. 把 `user_id`、`username` 放进上下文

### 7.2 RequirePermission

只负责：

1. 根据 `user_id` 查询角色与权限
2. 判断是否包含指定权限点
3. 无权限时返回统一 `403`

推荐使用方式：

```text
RequirePermission("article:create")
RequirePermission("article:update")
RequirePermission("category:create")
```

## 8. 接口规划

### 8.1 当前用户权限

- `GET /api/admin/rbac/me`

作用：

- 返回当前用户的角色和权限点
- 方便前端决定是否显示某些管理入口

### 8.2 角色管理

- `GET /api/admin/roles`
- `POST /api/admin/roles`

### 8.3 权限管理

- `GET /api/admin/permissions`
- `POST /api/admin/permissions`

### 8.4 用户绑定角色

- `PUT /api/admin/users/:user_id/roles`

### 8.5 角色绑定权限

- `PUT /api/admin/roles/:role_id/permissions`

## 9. 当前实现状态

当前代码已经完成：

1. `users` 表登录
2. `bcrypt` 密码校验
3. `super_admin` 与 `editor` 双角色 seed
4. `admin` 与 `editor1` 两个测试用户 seed
5. 文章、分类、标签后台接口的路由级权限接入
6. `401 / 403 / 200` 三层语义验证通过

## 10. 与当前项目的衔接方式

### 10.1 Session 不需要未来重做业务层

当前先使用本地 cookie store 维持 Session，后续如果切到 Redis store，业务 handler 与权限中间件不需要大幅改写，主要变化在底层 session store。

### 10.2 auth 和 rbac 分层清晰

- `auth` 管身份
- `rbac` 管授权

这样不会把登录逻辑、角色逻辑、权限逻辑全混到一个文件里。

### 10.3 现有后台接口已开始逐步接入

当前已经接入或应继续接入的顺序：

1. `GET /api/admin/articles`
2. `POST /api/admin/articles`
3. `PUT /api/admin/articles/:id`
4. `DELETE /api/admin/articles/:id`
5. `GET /api/admin/categories`
6. `POST /api/admin/categories`
7. `GET /api/admin/tags`
8. `POST /api/admin/tags`
9. 后续 `upload`

## 11. 本阶段成功标准

当这一阶段完成时，至少应满足：

1. 用户可以通过 Session 成功登录与退出
2. 路由能按权限点控制访问
3. 无登录返回 `401`
4. 已登录但无权限返回 `403`
5. 高权限与低权限用户都能稳定验证不同访问结果
6. RBAC 结构能自然承接后续上传和更多后台模块
