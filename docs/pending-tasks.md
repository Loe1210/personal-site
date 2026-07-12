# 个人小站待办事项与后续开发背景

## 0. 文档目的

本文档基于当前项目所有规划文档与开发日志（phase-01 至 phase-10）梳理而成，用于：

1. 明确当前项目已完成的边界
2. 列出规划中尚未完成的事项
3. 为后续阶段开发提供背景与参考

本文档会随项目推进持续更新，每完成一项应在此处标记完成状态。

---

## 1. 当前项目完成状态总览

### 1.1 已完成的核心闭环

| 模块 | 状态 | 说明 |
| --- | --- | --- |
| 项目骨架 | 已完成 | Hertz 单体 + thrift 驱动模型层 |
| 数据库持久化 | 已完成 | MySQL + GORM，含连接池配置 |
| 文章模块 | 已完成 | 公开读取 + 后台 CRUD，数据库分页，N+1 已优化 |
| 标签模块 | 已完成 | 独立 CRUD + 文章多对多关系 |
| 分类模块 | 已完成 | 独立 CRUD + 文章关联校验 |
| 认证模块 | 已完成 | Session 登录 + bcrypt 密码校验 |
| RBAC 权限校验 | 已完成 | 路由级权限中间件，401/403/200 三层语义 |
| 上传模块 | 已完成 | 本地存储 + 图片校验 + 与文章封面联动 |
| 前台页面 | 已完成 | 纯静态 SPA（首页/Blog 列表/详情/About） |
| 后台管理 | 已完成 | 纯静态 SPA（登录/仪表盘/文章/分类标签） |
| 前后端架构 | 已完成 | 从服务端模板渲染迁移到纯静态 SPA |
| 统一错误处理 | 已完成 | AppError + errno + response 统一封装 |

### 1.2 已完成的关键优化（phase-10）

- 数据库连接池配置（MaxIdleConns / MaxOpenConns / ConnMaxLifetime）
- 文章列表数据库级分页（Offset + Limit）
- 文章列表 N+1 查询修复（批量查询标签）
- 所有 Handler 统一错误处理
- 新增基于 ID 的文章查询接口 `GET /api/articles/id/:id`
- 删除默认测试用户 `editor1`
- 移除登录页密码提示
- 前端骨架屏、图片懒加载、错误 toast、标签筛选
- 静态资源缓存中间件
- 代码高亮语言包升级（支持全语言）

---

## 2. 待办事项分类

### 2.1 后端工程化完善（阶段5，完全未开始）

**来源**：`docs/personal-site-backend-plan.md` 第13节

这是后端规划中明确列出的"阶段5：工程化完善"，目前所有子项均未开始。

| 编号 | 事项 | 优先级 | 说明 |
| --- | --- | --- | --- |
| P-01 | Docker 化 | 中 | 容器化部署，便于环境一致 |
| P-02 | 环境配置拆分 | 中 | 当前依赖环境变量，计划引入配置文件（server/database/jwt/upload/site） |
| P-03 | 测试补齐 | 中 | 单元测试与接口测试 |
| P-04 | 上传能力替换为 MinIO | 低 | 当前本地存储可用，后续切对象存储 |
| P-05 | Redis 缓存接入 | 低 | 计划标注为"可选"，session store 可切 Redis |

**建议**：在 Kitex 迁移之前优先完成 Docker 化和配置拆分，为后续服务化打好基础。

---

### 2.2 RBAC 管理接口（完全未实现）

**来源**：`docs/personal-site-rbac-plan.md` 第8节

当前只完成了权限校验中间件和种子数据，以下管理接口均未实现：

| 编号 | 接口 | 用途 |
| --- | --- | --- |
| P-06 | `GET /api/admin/rbac/me` | 返回当前用户的角色和权限点 |
| P-07 | `GET /api/admin/roles` | 角色列表 |
| P-08 | `POST /api/admin/roles` | 角色创建 |
| P-09 | `GET /api/admin/permissions` | 权限点列表 |
| P-10 | `POST /api/admin/permissions` | 权限点创建 |
| P-11 | `PUT /api/admin/users/:user_id/roles` | 用户绑定角色 |
| P-12 | `PUT /api/admin/roles/:role_id/permissions` | 角色绑定权限 |

**说明**：`rbac.thrift` 第一版接口草案已起草并保留，可作为后续实现的契约基础。当前权限数据全部由启动时 seed 注入，无法通过后台动态管理。

**建议**：如果后续需要多用户协作，优先完成 P-06（前端可根据权限点决定是否显示管理入口）。

---

### 2.3 Kitex 微服务迁移（未开始，当前仅预留边界）

**来源**：`docs/personal-site-backend-plan.md` 第14节

当前项目已按领域拆分 service 边界（`ArticleService` / `AuthService` / `UploadService`），为后续 RPC 化预留了空间，但尚未开始实际迁移。

| 编号 | 事项 | 说明 |
| --- | --- | --- |
| P-13 | `article` 模块拆分为 RPC 服务 | 最容易拆出去的模块 |
| P-14 | `auth` 模块拆分为 RPC 服务 | 第二优先级 |
| P-15 | `upload` 模块拆分为 RPC 服务 | 第三优先级 |

**前提条件**：建议在完成 Docker 化（P-01）后再启动微服务拆分。

---

### 2.4 Eino AI 能力（未开始）

**来源**：`docs/personal-site-backend-plan.md` 第15节

Eino 定位为增强层，不是基础设施层，应在 Kitex 迁移之后接入。

| 编号 | 事项 | 说明 |
| --- | --- | --- |
| P-16 | 自动生成文章摘要 | 文章创建/更新时触发 |
| P-17 | 自动推荐标签 | 根据文章内容推荐 |
| P-18 | 根据文章生成学习卡片 | 知识沉淀 |
| P-19 | 站内问答 | 基于 Eino 的 RAG 能力 |

**建议**：这是最后阶段的能力，当前不优先。

---

### 2.5 功能遗留项

这些是开发过程中发现但未完成的功能性问题。

| 编号 | 事项 | 来源 | 说明 |
| --- | --- | --- | --- |
| P-20 | 上一页/下一篇文章导航 | UI 规范第8.3节 | phase-10 中移除了 `GetPrevNext` 函数，前端 `getAdjacentPosts` 返回空 |
| P-21 | About 页后台可编辑 | 后端规划第5.7节 | 当前为静态内容，计划后续转为后台可编辑配置项 |
| P-22 | 滑动黑线问题 | phase-10 第27节 | 尝试修复但未解决，用户反馈"还是有黑条"，可能与 CSS `transform` 或 `overflow` 有关 |
| P-23 | RefreshUserToken 接口 | IDL 规划第14.2节 | 计划标注"第一版可不做"，尚未实现 |
| P-24 | 限流中间件 | 后端规划第10.4节 | 计划标注"第一版可以不做"，尚未实现 |

**建议**：P-20 和 P-22 影响用户体验，可在下一阶段优先处理。P-21 视内容更新频率决定。

---

### 2.6 管理后台视觉对齐（需确认）

**来源**：`docs/personal-site-ui-redesign-plan.md` 第12节

phase-10 前端架构迁移后，新的纯静态 `static/admin/index.html` 是否已完全对齐 Terminal Gallery 设计系统需要确认。

| 编号 | 事项 | 说明 |
| --- | --- | --- |
| P-25 | 管理后台视觉对齐 Terminal Gallery | phase-09.1 曾做过 admin redesign，但 phase-10 删除了原 `templates/` 目录，admin 改为纯静态 SPA，视觉风格可能需要重新检查 |

**说明**：phase-10 中新建的 `static/admin/index.html` 和 `static/admin/css/admin.css` 是独立实现，是否继承 phase-09.1 的 Terminal Gallery 视觉需要核对。

---

## 3. 建议的后续开发顺序

基于以上分析，建议按以下顺序推进：

### 阶段 A：功能补全与体验修复（短期）

优先处理影响用户体验和功能完整性的遗留项：

1. **P-20** 上一页/下一篇文章导航
2. **P-22** 滑动黑线问题排查
3. **P-25** 管理后台视觉对齐确认

### 阶段 B：工程化基础（中期）

为后续服务化演进打基础：

4. **P-02** 环境配置拆分
5. **P-01** Docker 化
6. **P-03** 测试补齐

### 阶段 C：RBAC 管理（按需）

如果需要多用户协作：

7. **P-06** 当前用户权限查询接口
8. **P-07 ~ P-12** 角色与权限管理接口

### 阶段 D：微服务演进（长期）

9. **P-13 ~ P-15** Kitex 服务拆分
10. **P-04** 上传能力替换为 MinIO
11. **P-05** Redis 缓存接入

### 阶段 E：AI 增强（远期）

12. **P-16 ~ P-19** Eino 能力接入

---

## 4. 关键约束与注意事项

### 4.1 前端架构约束

phase-10 已完成前端架构迁移：
- 前端为纯静态 SPA，位于 `static/` 目录
- 后端只提供 API，不再做服务端模板渲染
- `biz/site/handler.go` 仅保留 Swagger UI
- 静态文件由 `main.go` 中的 `serveSPA` 处理器统一处理

后续前端开发应继续在 `static/` 目录下迭代，不要回退到服务端渲染。

### 4.2 认证架构约束

- 当前使用 Session 认证（非 JWT）
- 登录态通过 cookie store 维持
- 后续切 Redis session store 时，业务 handler 与权限中间件不需要大幅改写

### 4.3 IDL 约束

- 核心接口使用 thrift 定义
- 每个 thrift 有独立 `package` 和 `go_package`
- 生成目录按领域分包，避免互相覆盖
- `biz/model/*/custom.go` 存放手工补充的结构体（如 `GetCategoryRequest`），待 `hz` 工具链恢复后可由生成代码替换

### 4.4 编码规范

- 所有文档统一使用 UTF-8 无 BOM 编码
- phase-01 至 phase-09 曾出现 GBK 编码导致的乱码问题，已全部修复
- 后续新增文档必须使用 UTF-8 编码

---

## 5. 参考文档索引

| 文档 | 路径 | 用途 |
| --- | --- | --- |
| 后端开发规划 | `docs/personal-site-backend-plan.md` | 后端整体架构与阶段划分 |
| 前端开发规划 | `docs/personal-site-frontend-plan.md` | 前端页面范围与实现路线 |
| IDL 驱动方案 | `docs/personal-site-idl-plan.md` | thrift 拆分与接口契约 |
| RBAC 方案 | `docs/personal-site-rbac-plan.md` | 权限体系设计 |
| Admin 前端方案 | `docs/personal-site-admin-frontend-plan.md` | 管理后台页面规划 |
| UI 重设计方案 | `docs/personal-site-ui-redesign-plan.md` | Terminal Gallery 视觉方向 |
| UI 设计文档 | `docs/personal-site-ui-spec.md` | 页面级视觉规范 |
| 项目规范 | `docs/project-conventions.md` | 目录/Git/提交/日志规范 |
| 开发日志 | `docs/devlog/phase-01.md` ~ `phase-10.md` | 各阶段开发记录 |
| **本文档** | `docs/pending-tasks.md` | 待办事项与后续开发背景 |

---

## 6. 更新规则

- 每完成一项待办事项，应在对应编号后标记 `[已完成]` 并注明完成阶段
- 新发现的待办事项应追加到对应分类下
- 阶段推进后应更新"建议的后续开发顺序"
- 本文档不替代开发日志，具体实现细节仍记录在 `docs/devlog/` 中
