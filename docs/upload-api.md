# 本地文件上传 API

## 小文件上传

`POST /api/media/upload`

表单字段：

- `file`：必填，图片文件。
- `biz_type`：可选，业务类型，例如 `article_cover`。
- `sha256`：可选，文件 SHA-256；传入后服务端会校验内容 hash。

说明：该接口保留给封面图等小文件使用。服务端会校验文件类型、文件头和可选 hash。

## 分片上传流程

### 1. 初始化任务

`POST /api/media/uploads/init`

表单字段：

- `user_id`：必填，当前用户 ID。
- `file_name`：必填，原始文件名。
- `file_size`：必填，文件总大小。
- `chunk_size`：必填，分片大小。
- `content_type`：建议传入，文件 MIME 类型。
- `sha256`：可选，完整文件 SHA-256。
- `biz_type` / `biz_id`：可选，业务归属。

返回 `upload_id`、`chunk_count`、过期时间和当前状态。

### 2. 上传分片

`PUT /api/media/uploads/{upload_id}/chunks/{chunk_index}`

查询或表单字段：

- `user_id`：必填，必须是任务创建者。

请求体直接传当前分片内容。服务端会把分片流式写入 `static/uploads/tmp/{upload_id}`，并记录分片大小和 hash。

### 3. 查询/断点恢复

`GET /api/media/uploads/{upload_id}?user_id=1`

返回任务状态和已经上传的分片列表，前端可据此跳过已完成分片。

### 4. 完成上传

`POST /api/media/uploads/{upload_id}/complete`

表单字段：

- `user_id`：必填，必须是任务创建者。

服务端会按分片顺序合并到正式目录，校验完整文件 hash，图片会生成缩略图，然后在同一个事务里保存文件记录并把任务置为 `completed`。

### 5. 取消上传

`POST /api/media/uploads/{upload_id}/cancel`

表单字段：

- `user_id`：必填，必须是任务创建者。

## 状态机

- `uploading`：任务创建后、分片上传中。
- `completed`：完成合并并保存文件记录。
- `cancelled`：用户取消任务。
- `failed`：预留失败状态。

## 清理策略

media-service 启动后台回收任务，默认每 10 分钟扫描一次过期的 `uploading` 任务，删除对应临时目录和任务记录。
