# content-service 迁移说明

## 为什么 article/category/tag 放在一个内容域

文章、分类、标签共同描述博客内容。文章列表、详情、后台编辑都需要同时理解分类和标签关系，如果拆成多个服务，会让一次文章查询变成多次跨服务调用，反而增加一致性和排障成本。

本阶段把 `article`、`category`、`tag` 统一放入 `content-service`，让内容域自己拥有 `content_db` 和迁移脚本。其他服务只能通过 HTTP 或后续 Kitex RPC 访问内容能力，不能直接读取内容库表。

## 本阶段迁移范围

- 新增 `services/content-service` 独立服务目录。
- 应用层提供 `GetArticleByID`、`ListPublicArticles`、后台 CRUD 等最小边界。
- 公开文章详情统一按文章 `id` 查询。
- 新增 `content_db` 迁移脚本，包含 `articles`、`categories`、`tags`、`article_tags`。
- 新增 `idl/content/content.thrift`，作为后续 Kitex 契约基础。

## 当前保留的边界

这一阶段先建立服务骨架和可编译边界，原单体接口还没有被删除。后续接入 gateway 后，公开流量会通过 `gateway -> web-bff -> content-service` 进入内容域。

`content-service` 不处理登录态，也不直接读取 `auth-service` 数据库。后台权限由 gateway 或调用方通过 `auth-service` 校验后再访问内容服务。