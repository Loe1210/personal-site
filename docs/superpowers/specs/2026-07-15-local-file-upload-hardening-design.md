# 本地文件上传加固设计

**目标：** 在不引入对象存储和分发层的前提下，把现有 media-service 升级为支持分片上传、断点续传、秒传、临时区、图片处理、限流、超时和安全校验的本地文件服务。

**范围：** 仅覆盖文件上传链路与文件元数据治理，不改博客内容业务的查询/发布语义，不新增 CDN/下载分发层。

## 背景

当前上传接口会一次性读取完整文件到内存，适合封面图这类小文件，不适合大文件。为避免大文件和高并发上传拖垮服务，需要把上传从“单请求收文件”升级为“任务化、分片化、可恢复”的流程，同时保留小文件普通上传能力。

## 总体方案

采用“接入层 -> 上传层 -> 临时区存储 -> 图片处理 -> 元数据持久化 -> 正式文件落盘”的链路。

- 接入层负责统一鉴权、请求大小限制、基础限流和超时保护。
- 上传层负责创建任务、接收分片、查询进度、断点恢复、合并完成和取消任务。
- 存储层继续使用本地磁盘，但拆分临时区和正式区。
- 处理层只处理图片类文件，做压缩和缩略图生成。
- 元数据层记录文件、任务、分片、hash、所属用户/业务和版本号。
- 安全层负责文件头校验、hash 校验、类型校验和权限隔离。

现有普通上传接口保留，用于封面图等小文件；大文件走分片上传任务接口。

## 接口设计

### 1. 创建上传任务

`POST /api/media/uploads/init`

请求字段：
- `file_name`
- `file_size`
- `chunk_size`
- `content_type`
- `sha256`（可选，用于秒传和完整性校验）
- `biz_type`
- `biz_id`（可选）

返回字段：
- `upload_id`
- `chunk_size`
- `chunk_count`
- `expires_at`
- `status`
- `uploaded_chunks`

### 2. 上传分片

`PUT /api/media/uploads/:upload_id/chunks/:index`

请求体为二进制分片内容，服务端不整块读入内存，直接流式写入临时区。

返回字段：
- `upload_id`
- `chunk_index`
- `received`
- `status`

### 3. 查询上传状态

`GET /api/media/uploads/:upload_id`

返回字段：
- 任务状态
- 已上传分片列表
- 文件大小
- 已上传字节数
- 过期时间
- 最近错误原因

### 4. 完成上传

`POST /api/media/uploads/:upload_id/complete`

完成时执行：
- 分片完整性校验
- 文件头校验
- hash 校验
- 合并到正式目录
- 图片压缩和缩略图生成
- 写入文件元数据

### 5. 取消上传

`DELETE /api/media/uploads/:upload_id`

取消任务并清理临时区、任务记录和分片记录。

## 状态机

上传任务状态固定为：
- `pending`
- `uploading`
- `merging`
- `processing`
- `completed`
- `failed`
- `expired`
- `canceled`

状态变化只允许单向推进，失败后可通过断点恢复继续回到 `uploading`，但任务一旦 `completed` 就不可再写入。

## 数据设计

### upload_tasks

记录任务级状态和恢复信息：
- `upload_id`
- `user_id`
- `biz_type`
- `biz_id`
- `file_name`
- `file_size`
- `chunk_size`
- `chunk_count`
- `uploaded_chunks`
- `status`
- `sha256`
- `expires_at`
- `last_error`
- `version`
- `created_at`
- `updated_at`

### upload_chunks

记录分片级状态：
- `upload_id`
- `chunk_index`
- `chunk_size`
- `chunk_hash`
- `tmp_path`
- `status`
- `created_at`

### file_records

记录最终正式文件：
- `upload_id`
- `original_name`
- `storage_path`
- `public_url`
- `content_type`
- `size`
- `sha256`
- `biz_type`
- `owner_id`
- `version`
- `created_at`

## 本地存储布局

- 临时区：`static/uploads/tmp/{upload_id}/`
- 正式区：`static/uploads/images/YYYYMMDD/`
- 缩略图：`static/uploads/images/YYYYMMDD/thumbs/`

临时区中的文件只服务于未完成任务，任务完成后移动到正式区；任务过期后统一清理。

## 图片处理

仅对图片类型执行：
- 压缩
- 生成缩略图
- 可选尺寸标准化

处理应放在上传完成后的独立步骤中，不阻塞分片接收。

## 可靠性设计

- 幂等：同一 `upload_id + chunk_index` 重复提交只保留一次成功结果。
- 断点恢复：前端先查询已上传分片，再补传缺失分片。
- 失败重试：分片上传失败后允许重试同一分片。
- 任务回收：过期任务及其临时文件由定时任务清理。
- 超时保护：单分片上传、任务完成和图片处理都设置超时。

## 安全设计

- 接入层统一鉴权。
- 上传接口限制请求体大小、单任务大小和并发数。
- 仅允许白名单图片类型。
- 校验文件头，不信任扩展名和 `Content-Type`。
- 校验总文件 hash，避免分片错位或篡改。
- 任务只能由创建者访问和恢复。

## 现有系统的落点

- `services/media-service` 承担上传任务、临时区和文件元数据。
- `frontend/nginx` 继续负责本地静态文件访问。
- `services/gateway` 可继续做统一入口和限流前置。
- `static/admin` 后台继续调用媒体服务上传接口，封面图仍可走普通小文件上传。

## 迁移原则

- 先保留现有 `/api/media/upload`，不破坏封面图上传。
- 先实现大文件分片链路，再逐步把后台大文件上传入口切过去。
- 只改 media-service 内部结构，不引入 OSS/CDN。

## 验证方式

- 单元测试覆盖任务状态机、幂等、hash 校验和文件类型校验。
- 集成测试覆盖初始化、分片上传、断点恢复、完成合并、取消任务和过期清理。
- 手工验证小文件仍可通过旧接口上传，大文件通过分片接口完成。
- 验证临时区和正式区都能被正确清理和访问。
