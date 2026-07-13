# media-service 迁移说明

## 为什么早拆上传服务

上传能力天然适合先从单体中拆出来：它和文章、分类、标签的业务关系较弱，但会引入文件存储、访问路径、大小限制、类型校验和元数据表。提前拆成 `media-service` 后，后续 `content-service` 只需要保存图片 URL 或文件 ID，不需要关心文件落在哪里。

## 本阶段迁移范围

- 新增独立服务目录 `services/media-service`。
- 提供 HTTP `/upload` 入口，用于接收表单字段 `file` 和可选 `biz_type`。
- 提供 HTTP `/files/:id` 入口，用于按文件 `id` 查询元数据。
- 新增 `files` 表作为 media-service 自己的数据库 schema。
- 新增 `idl/media/media.thrift`，先定义元数据查询 RPC 契约。

## 当前为什么先用本地存储

第一阶段目标是把服务边界跑通，而不是一次性引入对象存储。当前本地存储实现位于 `internal/infra/storage/local.go`，它只负责把文件写入配置中的 `upload.root_dir`，并返回可访问路径。

这样做的好处是迁移成本低、验证链路短：HTTP 上传、应用层服务、MySQL 元数据和访问 URL 可以先闭环。等 gateway、Nacos、K8s 和统一配置完成后，再把存储实现替换为 MinIO 或其他对象存储。

## 后续切换 MinIO 的替换点

`media-service` 应用层只依赖 `Storage` 接口：

```go
type Storage interface {
	Save(name string, content []byte) (string, error)
}
```

后续切换 MinIO 时，只需要新增 `internal/infra/storage/minio.go` 并在 `cmd/main.go` 根据配置选择实现。应用层、HTTP handler、MySQL repository 不需要改动。

## 和其他服务的边界

- `media-service` 拥有自己的 `media_db.files` 表。
- `content-service` 后续只保存文件 URL 或文件 ID，不直接读取 `files` 表。
- 文件上传继续使用 HTTP 更合适；大文件内容不通过 Kitex RPC 传输。
- 元数据查询可以通过 Kitex RPC 暴露给 `web-bff` 或后台聚合层。
