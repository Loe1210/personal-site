# Phase 05

## 目标

完成后台认证与授权主线从原型态到可运行第一版的落地：用 `session` 取代原来的 JWT 主链，用 `users` 表承接真实登录用户，并让 RBAC 最小闭环真正接入后台路由。

## 完成内容

1. 为接口错误处理补充结构化业务错误模型
2. 新增统一的成功/失败响应写出方法
3. 将认证主链从 JWT 切换为 Session
4. 在 `main.go` 接入 session 中间件与 cookie store
5. 将 `auth.thrift` 收口为当前 `User / UserLogin / GetCurrentUser` 结构
6. 登录接口改为写入 session，不再返回 token 与过期时间
7. `/api/admin/me` 改为基于 session 获取当前用户
8. 新增 `/api/admin/logout`，支持清理当前登录 session
9. 建立 `users` 表并用 `bcrypt` 存储密码哈希
10. 启动时自动 seed 默认用户 `admin`
11. 启动时自动 seed 低权限测试用户 `editor1`
12. 建立 `roles / permissions / user_roles / role_permissions` 四张 RBAC 关系表
13. 启动时自动 seed `super_admin` 与 `editor` 两个角色
14. 启动时自动 seed 基础权限点，并完成角色权限绑定
15. 为 `admin` 绑定 `super_admin`，为 `editor1` 绑定 `editor`
16. 实现 `RequirePermission(permission)` 路由级权限中间件
17. 为文章后台 CRUD 接入权限校验
18. 为分类后台列表/创建接入权限校验
19. 为标签后台列表/创建接入权限校验
20. 验证通过 `401 / 403 / 200` 三层语义
21. 起草并保留 `rbac.thrift` 第一版接口草案，作为后续 RBAC 管理接口演进基础

## 设计决策

### 1. 认证主链从 JWT 切换为 Session

项目后续路线以 Session 为主，因此本阶段不再继续扩展 JWT，而是直接把登录态切换为：

- 登录成功后写入 session
- 受保护接口通过 session 判断是否登录
- 用户退出时清理 session

这样可以让后续 Redis session 化和后台管理场景自然衔接。

### 2. 先用 `users` 承接真实用户，再做权限控制

RBAC 只有绑定在真实用户上才有意义，因此本阶段先完成：

- `users` 表
- 数据库登录校验
- `bcrypt` 密码比对

然后再把 `roles / permissions / user_roles / role_permissions` 接上。

### 3. 认证和授权分层保持清晰

当前代码里已经明确拆成两层：

- `AuthMiddleware()` 只负责确认是否已登录
- `RequirePermission()` 只负责确认是否有权限

这样后面继续扩展后台接口时，不需要把“登录”和“权限”逻辑混到一起。

### 4. 先做路由级权限，再考虑更细的数据级权限

当前阶段优先保证后台接口入口受到控制，因此先把权限挂在路由上。像“是否能修改某一篇特定文章”这类数据级权限，留到后续真正有业务需要时再补。

## 遇到的问题

### 1. Windows 环境中的 thrift 文件 BOM 导致 `hz model` 解析失败

`idl/auth.thrift` 一度因为 UTF-8 BOM 导致 `line 1 symbol 1` 解析错误。最终通过无 BOM 的 UTF-8 重写方式解决，并恢复了 auth model 的正常生成。

### 2. Session 化之后，接口模型和 Swagger 需要同步收口

认证从 JWT 切到 Session 后，登录返回不应继续保留 `token / expires_at`。因此本阶段同步调整了 `auth.thrift`、model 和 Swagger 生成结果，让接口语义与实际实现保持一致。

### 3. 中间件未中断链路会造成双响应

在早期 session 认证接入中，未登录分支如果只写错误响应但不 `Abort()`，后续 handler 仍会继续执行，导致响应体中出现两段 JSON。这个问题在 session 认证中间件中已修正。

## 当前结果

当前项目已经具备以下能力：

- 用数据库用户 + `bcrypt` 实现后台登录
- 用 Session 维持后台登录态
- 支持当前用户查询与退出登录
- 用角色与权限关系控制后台接口访问
- 同时拥有高权限测试账号与低权限测试账号，能够稳定验证 `401 / 403 / 200` 场景

当前可用的典型账号：

- `admin`：绑定 `super_admin`
- `editor1`：绑定 `editor`

## 下一步

当前阶段已经完成第一版 RBAC 最小闭环。下一步建议：

1. 清理剩余旧 JWT 文档表述，避免后续阅读混淆
2. 视需要继续为更多后台接口挂权限
3. 进入 `upload` 模块，为文章封面图和站点资源提供上传能力
4. 后续再补 RBAC 管理接口与更细粒度的权限控制
