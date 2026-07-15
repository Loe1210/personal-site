# 本地文件上传加固 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在保留本地磁盘存储的前提下，把 media-service 升级为支持分片上传、断点续传、秒传、临时区、图片压缩/缩略图、限流、超时、幂等和安全校验的文件服务。

**Architecture:** 接入层负责鉴权、请求体限制和限流；上传层负责任务状态机、分片接收、断点恢复和完成合并；存储层继续用本地盘，但拆分临时区和正式区；处理层只处理图片压缩和缩略图；元数据层保存任务、分片、文件、hash、所属用户/业务和版本号。

**Tech Stack:** Go, Hertz, GORM, MySQL, 本地文件系统, Docker Compose, 原生前端 JavaScript

## Global Constraints

- 保留现有 `/api/media/upload`，封面图等小文件继续可用。
- 不引入 OSS、MinIO 或 CDN 分发层。
- 上传接口必须支持本地盘临时区和正式区分离。
- 分片接收不得一次性把大文件读入内存。
- 任务只能由创建者访问、恢复和取消。
- 图片类文件才做压缩和缩略图生成。

---

### Task 1: 上传任务与元数据模型

**Files:**
- Create: `services/media-service/internal/model/upload_task.go`
- Create: `services/media-service/internal/model/upload_chunk.go`
- Modify: `services/media-service/internal/model/file.go`
- Create: `services/media-service/internal/dal/db/upload_task_repository.go`
- Create: `services/media-service/internal/dal/db/upload_chunk_repository.go`
- Modify: `services/media-service/internal/dal/db/migrate.go`
- Create: `services/media-service/migrations/002_upload_tasks.sql`

**Interfaces:**
- Consumes: `upload_id`, `file_name`, `file_size`, `chunk_size`, `content_type`, `sha256`, `biz_type`, `biz_id`, `user_id`
- Produces: `upload_tasks`, `upload_chunks`, `file_records`

- [ ] **Step 1: Write the failing test**

```go
func TestUploadTaskRepositoryStoresStateAndChunks(t *testing.T) {
    // create task, save one chunk, reload task, assert uploaded_chunks and status
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/media-service/internal/dal/db -run TestUploadTaskRepositoryStoresStateAndChunks -count=1`
Expected: FAIL because repositories and models do not exist yet.

- [ ] **Step 3: Write minimal implementation**

Add models and repository methods:

```go
type UploadTask struct {
    UploadID       string
    UserID         int64
    BizType        string
    BizID          string
    FileName       string
    FileSize       int64
    ChunkSize      int64
    ChunkCount     int
    UploadedChunks  string
    Status         string
    Sha256         string
    ExpiresAt      time.Time
    LastError      string
    Version        int64
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./services/media-service/internal/dal/db ./services/media-service/internal/model`
Expected: PASS.

---

### Task 2: Upload task service and state machine

**Files:**
- Create: `services/media-service/internal/service/upload_task_service.go`
- Create: `services/media-service/internal/service/upload_task_service_test.go`
- Modify: `services/media-service/cmd/main.go`
- Modify: `services/media-service/cmd/router.go`
- Modify: `services/media-service/biz/router.go`
- Create: `services/media-service/biz/upload/task_handler.go`
- Create: `services/media-service/biz/upload/task_route.go`

**Interfaces:**
- Consumes: repository methods from Task 1
- Produces: `InitUpload`, `GetUpload`, `CancelUpload`, `CompleteUpload`

- [ ] **Step 1: Write the failing test**

```go
func TestInitUploadRejectsTooLargeFile(t *testing.T) {
    // init with file_size larger than configured limit, expect error
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/media-service/internal/service -run TestInitUploadRejectsTooLargeFile -count=1`
Expected: FAIL because the service does not yet exist.

- [ ] **Step 3: Write minimal implementation**

Implement task state transitions:

```go
func (s *UploadTaskService) InitUpload(ctx context.Context, in InitInput) (*UploadTask, error)
func (s *UploadTaskService) GetUpload(ctx context.Context, uploadID string) (*UploadTask, []UploadedChunk, error)
func (s *UploadTaskService) CancelUpload(ctx context.Context, uploadID string) error
func (s *UploadTaskService) CompleteUpload(ctx context.Context, uploadID string) (*FileRecord, error)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./services/media-service/internal/service -run TestInitUploadRejectsTooLargeFile -count=1`
Expected: PASS.

---

### Task 3: 分片接收、临时区和断点恢复

**Files:**
- Create: `services/media-service/internal/dal/storage/tmp_storage.go`
- Create: `services/media-service/internal/dal/storage/tmp_storage_test.go`
- Modify: `services/media-service/internal/service/media_service.go`
- Create: `services/media-service/internal/service/chunk_service.go`
- Create: `services/media-service/internal/service/chunk_service_test.go`
- Modify: `services/media-service/biz/upload/handler.go`
- Modify: `services/media-service/biz/upload/route.go`

**Interfaces:**
- Consumes: upload task repository, chunk repository, tmp storage
- Produces: streamed chunk writes, uploaded chunk index tracking, merge-ready tmp files

- [ ] **Step 1: Write the failing test**

```go
func TestChunkServiceWritesChunkToTmpPath(t *testing.T) {
    // upload one chunk, assert temp file exists and metadata is recorded
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/media-service/internal/service -run TestChunkServiceWritesChunkToTmpPath -count=1`
Expected: FAIL because chunk service and tmp storage are missing.

- [ ] **Step 3: Write minimal implementation**

Stream chunk bytes directly to `static/uploads/tmp/{upload_id}/chunk_000001.part` and mark the chunk as uploaded only after the write succeeds.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./services/media-service/internal/service ./services/media-service/internal/dal/storage`
Expected: PASS.

---

### Task 4: 图片压缩、缩略图和完成合并

**Files:**
- Create: `services/media-service/internal/service/image_processor.go`
- Create: `services/media-service/internal/service/image_processor_test.go`
- Modify: `services/media-service/internal/service/upload_task_service.go`
- Modify: `services/media-service/internal/dal/storage/local.go`
- Create: `services/media-service/internal/service/merge_service.go`
- Create: `services/media-service/internal/service/merge_service_test.go`

**Interfaces:**
- Consumes: completed tmp files, content type, sha256
- Produces: final file URL, thumbnail URL, merged file record

- [ ] **Step 1: Write the failing test**

```go
func TestImageProcessorCreatesThumbnail(t *testing.T) {
    // input a png, expect thumbnail file generated under thumbs/
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/media-service/internal/service -run TestImageProcessorCreatesThumbnail -count=1`
Expected: FAIL because the processor does not exist yet.

- [ ] **Step 3: Write minimal implementation**

Compress images after merge, generate thumbnail beside the final asset, and keep the original file path as the canonical public URL.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./services/media-service/internal/service ./services/media-service/internal/dal/storage`
Expected: PASS.

---

### Task 5: 接入层限流、请求体限制、超时和安全校验

**Files:**
- Create: `services/gateway/internal/middleware/upload_guard.go`
- Create: `services/gateway/internal/middleware/upload_guard_test.go`
- Modify: `services/gateway/internal/router/router.go`
- Modify: `services/media-service/biz/upload/handler.go`
- Create: `services/media-service/internal/service/upload_security.go`
- Create: `services/media-service/internal/service/upload_security_test.go`

**Interfaces:**
- Consumes: auth context, request size, mime type, file header, sha256
- Produces: rejected oversized requests, rate-limited requests, validated uploads

- [ ] **Step 1: Write the failing test**

```go
func TestUploadSecurityRejectsWrongMagicBytes(t *testing.T) {
    // pretend image type but wrong header, expect rejection
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/media-service/internal/service -run TestUploadSecurityRejectsWrongMagicBytes -count=1`
Expected: FAIL because security checks are not implemented.

- [ ] **Step 3: Write minimal implementation**

Add request size limits, per-user/per-IP rate limiting, file header validation, hash validation, and upload timeout checks.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./services/media-service/internal/service ./services/gateway/internal/middleware`
Expected: PASS.

---

### Task 6: 临时任务回收与断点恢复

**Files:**
- Create: `services/media-service/internal/service/upload_reaper.go`
- Create: `services/media-service/internal/service/upload_reaper_test.go`
- Modify: `services/media-service/cmd/main.go`

**Interfaces:**
- Consumes: upload task repository, tmp storage
- Produces: expired task cleanup, orphan chunk cleanup

- [ ] **Step 1: Write the failing test**

```go
func TestUploadReaperDeletesExpiredTmpFiles(t *testing.T) {
    // create expired task and tmp files, run reaper, expect cleanup
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/media-service/internal/service -run TestUploadReaperDeletesExpiredTmpFiles -count=1`
Expected: FAIL because the reaper does not exist yet.

- [ ] **Step 3: Write minimal implementation**

Run a background ticker that removes expired upload tasks and their tmp directories.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./services/media-service/internal/service`
Expected: PASS.

---

### Task 7: 文档与运行验证

**Files:**
- Modify: `services/media-service/configs/config.yaml`
- Modify: `configs/config.yaml`
- Modify: `deploy/docker/compose.yaml`
- Create: `docs/upload-api.md`

**Interfaces:**
- Consumes: new upload endpoints and config fields
- Produces: runnable local setup, tmp volume persistence, and usage notes

- [ ] **Step 1: Update config defaults for chunk uploads, temp TTL, and concurrency limits**

Update these config keys in `services/media-service/configs/config.yaml` and `configs/config.yaml`:

```yaml
upload:
  root_dir: "static/uploads/images"
  public_base_path: "/static/uploads/images"
  tmp_root_dir: "static/uploads/tmp"
  chunk_size_mb: 4
  max_upload_size_mb: 512
  max_concurrent_uploads: 10
  max_concurrent_chunks: 3
  upload_ttl_hour: 24
  cleanup_interval_min: 10
```

- [ ] **Step 2: Add and mount a persistent tmp volume for media-service and a public image volume for frontend**

Mount `upload-tmp-data:/app/static/uploads/tmp` in `media-service` and `upload-images-data:/usr/share/nginx/html/static/uploads/images` in `frontend`.

- [ ] **Step 3: Rebuild containers**

Run: `docker compose -f deploy/docker/compose.yaml up -d --build media-service gateway frontend`
Expected: media-service, gateway, and frontend all restart cleanly, and the shared upload volumes mount without errors.

- [ ] **Step 4: Verify the old small-file upload still works and a large-file chunk flow can start**

Run: `go test ./services/media-service/... ./services/gateway/...`
Expected: PASS.

- [ ] **Step 5: Add usage notes for init/chunk/complete/cancel**

Document the request shapes, task state machine, chunk retry sequence, and cleanup behavior.