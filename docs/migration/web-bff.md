# web-bff 迁移说明

## 为什么需要 web-bff

前端页面不应该直接拼多个后端服务接口。页面需要的是面向展示的 DTO，而不是领域服务的数据库形状。`web-bff` 负责把内容服务、媒体服务、认证上下文等结果组装成前端更容易消费的页面数据。

这样前端可以保持稳定：内容域内部表结构、RPC 契约或服务拆分变化时，只要 `web-bff` 输出的页面 DTO 不变，页面代码就不需要频繁跟着改。

## 当前 web-bff 只做什么

- 提供 `/blog/articles/:id` 页面聚合入口。
- 调用 `content-service` 获取文章详情。
- 输出 `ArticlePageDTO`，把文章详情包成页面数据。

## 当前 web-bff 不做什么

- 不拥有文章、分类、标签的领域真相。
- 不直接访问 `content_db`。
- 不做后台权限判定。
- 不承担文件存储能力。

后续任务 5 接入 gateway、Nacos 和 OpenTelemetry 后，`web-bff` 会从固定 HTTP 地址迁移到服务发现或 Kitex 客户端调用。