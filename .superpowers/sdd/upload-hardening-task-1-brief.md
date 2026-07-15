# Task 1 Brief

Source plan: docs/superpowers/plans/2026-07-15-local-file-upload-hardening.md

Global Constraints:
- 保留现有 `/api/media/upload`，封面图等小文件继续可用。
- 不引入 OSS、MinIO 或 CDN 分发层。
- 上传接口必须支持本地盘临时区和正式区分离。
- 分片接收不得一次性把大文件读入内存。
- 任务只能由创建者访问、恢复和取消。
- 图片类文件才做压缩和缩略图生成。

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

