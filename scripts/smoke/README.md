# 微服务 Smoke 验证

运行前先启动本地微服务：

```bash
make micro-up
```

然后执行：

```bash
make micro-smoke
```

脚本会验证：

- gateway `/healthz` 可访问。
- auth-service 未登录访问 `/me` 返回 401。
- content-service 文章列表接口可访问。
- auth-service 登录后能通过 session cookie 访问 `/me`。

如果任一步失败，脚本会以非 0 状态退出。