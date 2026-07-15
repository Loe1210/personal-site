# Upload Hardening Task 2 Review Package

- Base: `470d4b0`
- Head: `8623e92`
- Task: Upload task service and state machine

## Commits
8623e92 docs: record upload hardening task 2 report
2d69978 feat(media-service): add upload task service

## Diff Stat
```text
 .superpowers/sdd/upload-hardening-task-2-brief.md  |  58 ++++++++++
 .superpowers/sdd/upload-hardening-task-2-report.md |  22 ++++
 services/media-service/biz/router.go               |   5 +-
 services/media-service/biz/upload/task_handler.go  | 107 ++++++++++++++++++
 services/media-service/biz/upload/task_route.go    |  15 +++
 services/media-service/cmd/main.go                 |   6 +-
 services/media-service/cmd/router.go               |   6 +-
 .../internal/service/upload_task_service.go        | 123 +++++++++++++++++++++
 .../internal/service/upload_task_service_test.go   |  26 +++++
 9 files changed, 361 insertions(+), 7 deletions(-)
```

## Full Diff
```diff
diff --git a/.superpowers/sdd/upload-hardening-task-2-brief.md b/.superpowers/sdd/upload-hardening-task-2-brief.md
new file mode 100644
index 0000000..9bf1a2d
--- /dev/null
+++ b/.superpowers/sdd/upload-hardening-task-2-brief.md
@@ -0,0 +1,58 @@
+# Task 2 Brief - Upload task service and state machine
+
+You are implementing Task 2 of the local file upload hardening plan in `docs/superpowers/plans/2026-07-15-local-file-upload-hardening.md`.
+
+## Where this task fits
+Task 1 already added upload task/chunk persistence, migrations, and repository tests. This task must build the service/state-machine layer on top of those repositories so later chunk upload and merge work has a real API contract.
+
+## Requirements
+- Use the existing upload task/chunk repositories from Task 1.
+- Keep local disk storage only. Do not introduce OSS, MinIO, or CDN distribution.
+- Preserve the existing small-file upload flow and `/api/media/upload`.
+- The task state machine must support create, query, cancel, and complete.
+- A task is only accessible by its creator.
+- Follow the plan¡¯s intended interfaces:
+  - `InitUpload`
+  - `GetUpload`
+  - `CancelUpload`
+  - `CompleteUpload`
+
+## Files in scope
+- Create: `services/media-service/internal/service/upload_task_service.go`
+- Create: `services/media-service/internal/service/upload_task_service_test.go`
+- Modify: `services/media-service/cmd/main.go`
+- Modify: `services/media-service/cmd/router.go`
+- Modify: `services/media-service/biz/router.go`
+- Create: `services/media-service/biz/upload/task_handler.go`
+- Create: `services/media-service/biz/upload/task_route.go`
+
+## Context you should assume
+- `services/media-service/internal/model/upload_task.go` and `upload_chunk.go` already exist.
+- `services/media-service/internal/dal/db/upload_task_repository.go` and `upload_chunk_repository.go` already exist.
+- The database migration and test sqlite dependency are already in place.
+- Do not remove or rewrite unrelated changes made by other tasks.
+
+## Test-first target
+Start with the failing test from the plan:
+
+```go
+func TestInitUploadRejectsTooLargeFile(t *testing.T) {
+    // init with file_size larger than configured limit, expect error
+}
+```
+
+Then make it pass with the smallest change set that preserves the plan.
+
+## Report file
+Write your full report to:
+`.superpowers/sdd/upload-hardening-task-2-report.md`
+
+Report status should include:
+- what you implemented
+- tests run and results
+- any concerns or follow-up notes
+
+## Important
+- You are not alone in the codebase. Do not revert or overwrite edits made by others.
+- Keep your work tightly scoped to Task 2.
+- Commit your changes when done and include the commit hash in the report.
diff --git a/.superpowers/sdd/upload-hardening-task-2-report.md b/.superpowers/sdd/upload-hardening-task-2-report.md
new file mode 100644
index 0000000..9a40465
--- /dev/null
+++ b/.superpowers/sdd/upload-hardening-task-2-report.md
@@ -0,0 +1,22 @@
+# Task 2 Report - Upload task service and state machine
+
+Status: DONE
+
+## Completed
+- Added `UploadTaskService` with `InitUpload`, `GetUpload`, `CancelUpload`, and `CompleteUpload`.
+- Added upload task sizing/default chunk calculation and 24-hour task expiration default.
+- Added HTTP task endpoints under `/upload/tasks/*`.
+- Wired upload task repositories into `media-service` startup and route registration.
+- Preserved the existing `/upload` small-file flow.
+- Tightened the size-limit test so it validates the intended failure path.
+
+## Verification
+- `go test ./services/media-service/internal/service -run TestInitUploadRejectsTooLargeFile -count=1`
+- `go test ./services/media-service/...`
+
+## Notes
+- `CompleteUpload` currently only transitions status to `completed`; actual chunk merge and file record creation remain for later tasks.
+- Task ownership is enforced through repository lookups using `upload_id` plus `user_id`.
+
+## Commit
+- 2d69978
diff --git a/services/media-service/biz/router.go b/services/media-service/biz/router.go
index d8d91d3..eeed41f 100644
--- a/services/media-service/biz/router.go
+++ b/services/media-service/biz/router.go
@@ -1,12 +1,13 @@
 package biz
 
 import (
 	"github.com/cloudwego/hertz/pkg/app/server"
 
 	"github.com/Loe1210/personal-site/services/media-service/biz/upload"
 	"github.com/Loe1210/personal-site/services/media-service/internal/service"
 )
 
-func RegisterRoutes(hertz *server.Hertz, media *service.Service) {
+func RegisterRoutes(hertz *server.Hertz, media *service.Service, uploadTasks *service.UploadTaskService) {
 	upload.RegisterRoutes(hertz, media)
-}
\ No newline at end of file
+	upload.RegisterTaskRoutes(hertz, uploadTasks)
+}
diff --git a/services/media-service/biz/upload/task_handler.go b/services/media-service/biz/upload/task_handler.go
new file mode 100644
index 0000000..ec443e9
--- /dev/null
+++ b/services/media-service/biz/upload/task_handler.go
@@ -0,0 +1,107 @@
+package upload
+
+import (
+	"context"
+	"strconv"
+
+	"github.com/cloudwego/hertz/pkg/app"
+	"github.com/cloudwego/hertz/pkg/protocol/consts"
+
+	"github.com/Loe1210/personal-site/services/media-service/internal/service"
+)
+
+type TaskHandler struct {
+	service *service.UploadTaskService
+}
+
+func NewTaskHandler(service *service.UploadTaskService) *TaskHandler {
+	return &TaskHandler{service: service}
+}
+
+func (h *TaskHandler) InitUpload(ctx context.Context, c *app.RequestContext) {
+	userID, err := parseUploadUserID(c)
+	if err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20020, "message": err.Error()})
+		return
+	}
+	fileSize, err := parseFormInt64(c, "file_size")
+	if err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20021, "message": "invalid file size"})
+		return
+	}
+	chunkSize, err := parseFormInt64(c, "chunk_size")
+	if err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20022, "message": "invalid chunk size"})
+		return
+	}
+	task, err := h.service.InitUpload(ctx, service.InitInput{
+		UserID:      userID,
+		FileName:    c.PostForm("file_name"),
+		FileSize:    fileSize,
+		ContentType: c.PostForm("content_type"),
+		BizType:     c.PostForm("biz_type"),
+		BizID:       c.PostForm("biz_id"),
+		Sha256:      c.PostForm("sha256"),
+		ChunkSize:   chunkSize,
+	})
+	if err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20023, "message": err.Error()})
+		return
+	}
+	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": task})
+}
+
+func (h *TaskHandler) GetUpload(ctx context.Context, c *app.RequestContext) {
+	userID, err := parseUploadUserID(c)
+	if err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20024, "message": err.Error()})
+		return
+	}
+	task, chunks, err := h.service.GetUpload(ctx, c.Param("upload_id"), userID)
+	if err != nil {
+		c.JSON(consts.StatusNotFound, map[string]any{"code": 20025, "message": err.Error()})
+		return
+	}
+	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": map[string]any{"task": task, "chunks": chunks}})
+}
+
+func (h *TaskHandler) CancelUpload(ctx context.Context, c *app.RequestContext) {
+	userID, err := parseUploadUserID(c)
+	if err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20026, "message": err.Error()})
+		return
+	}
+	if err := h.service.CancelUpload(ctx, c.Param("upload_id"), userID); err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20027, "message": err.Error()})
+		return
+	}
+	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success"})
+}
+
+func (h *TaskHandler) CompleteUpload(ctx context.Context, c *app.RequestContext) {
+	userID, err := parseUploadUserID(c)
+	if err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20028, "message": err.Error()})
+		return
+	}
+	if err := h.service.CompleteUpload(ctx, c.Param("upload_id"), userID); err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20029, "message": err.Error()})
+		return
+	}
+	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success"})
+}
+
+func parseUploadUserID(c *app.RequestContext) (int64, error) {
+	return parseFormInt64(c, "user_id")
+}
+
+func parseFormInt64(c *app.RequestContext, key string) (int64, error) {
+	value := c.PostForm(key)
+	if value == "" {
+		value = c.Query(key)
+	}
+	if value == "" {
+		return 0, strconv.ErrSyntax
+	}
+	return strconv.ParseInt(value, 10, 64)
+}
diff --git a/services/media-service/biz/upload/task_route.go b/services/media-service/biz/upload/task_route.go
new file mode 100644
index 0000000..aa22c1c
--- /dev/null
+++ b/services/media-service/biz/upload/task_route.go
@@ -0,0 +1,15 @@
+package upload
+
+import (
+	"github.com/cloudwego/hertz/pkg/app/server"
+
+	"github.com/Loe1210/personal-site/services/media-service/internal/service"
+)
+
+func RegisterTaskRoutes(hertz *server.Hertz, uploadTasks *service.UploadTaskService) {
+	handler := NewTaskHandler(uploadTasks)
+	hertz.POST("/upload/tasks/init", handler.InitUpload)
+	hertz.GET("/upload/tasks/:upload_id", handler.GetUpload)
+	hertz.POST("/upload/tasks/:upload_id/cancel", handler.CancelUpload)
+	hertz.POST("/upload/tasks/:upload_id/complete", handler.CompleteUpload)
+}
diff --git a/services/media-service/cmd/main.go b/services/media-service/cmd/main.go
index be3d47d..bfa77ef 100644
--- a/services/media-service/cmd/main.go
+++ b/services/media-service/cmd/main.go
@@ -29,16 +29,18 @@ func main() {
 	}
 	defer shutdown(ctx)
 	database, err := db.Open(cfg.MySQL)
 	if err != nil {
 		log.Fatal(err)
 	}
 	if err := db.Migrate(database); err != nil {
 		log.Fatal(err)
 	}
 	store := storage.NewLocalStorage(cfg.Upload.RootDir, cfg.Upload.PublicBasePath)
-	media := service.NewMediaService(store, db.NewFileRepository(database))
+	fileRepo := db.NewFileRepository(database)
+	uploadTasks := service.NewUploadTaskService(&cfg.Upload, db.NewUploadTaskRepository(database), db.NewUploadChunkRepository(database))
+	media := service.NewMediaService(store, fileRepo)
 	startMediaRPCServer(cfg.RPC.Port, kitexmediahandler.NewHandler(media))
-	h := newRouter(media, configs.GetServerAddr())
+	h := newRouter(media, uploadTasks, configs.GetServerAddr())
 	log.Printf("media-service listening on %s", configs.GetServerAddr())
 	h.Spin()
 }
diff --git a/services/media-service/cmd/router.go b/services/media-service/cmd/router.go
index 42bcdda..c7321a1 100644
--- a/services/media-service/cmd/router.go
+++ b/services/media-service/cmd/router.go
@@ -1,14 +1,14 @@
 package main
 
 import (
 	"github.com/cloudwego/hertz/pkg/app/server"
 
 	"github.com/Loe1210/personal-site/services/media-service/biz"
 	"github.com/Loe1210/personal-site/services/media-service/internal/service"
 )
 
-func newRouter(media *service.Service, addr string) *server.Hertz {
+func newRouter(media *service.Service, uploadTasks *service.UploadTaskService, addr string) *server.Hertz {
 	h := server.Default(server.WithHostPorts(addr))
-	biz.RegisterRoutes(h, media)
+	biz.RegisterRoutes(h, media, uploadTasks)
 	return h
-}
\ No newline at end of file
+}
diff --git a/services/media-service/internal/service/upload_task_service.go b/services/media-service/internal/service/upload_task_service.go
new file mode 100644
index 0000000..bc3a906
--- /dev/null
+++ b/services/media-service/internal/service/upload_task_service.go
@@ -0,0 +1,123 @@
+package service
+
+import (
+	"context"
+	"errors"
+	"math"
+	"time"
+
+	"github.com/Loe1210/personal-site/configs"
+	db "github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
+	"github.com/Loe1210/personal-site/services/media-service/internal/model"
+	"github.com/google/uuid"
+)
+
+const defaultChunkSizeBytes int64 = 4 * 1024 * 1024
+
+type InitInput struct {
+	UserID      int64
+	FileName    string
+	FileSize    int64
+	ContentType string
+	BizType     string
+	BizID       string
+	Sha256      string
+	ChunkSize   int64
+}
+
+type UploadTaskService struct {
+	cfg         *configs.UploadConfig
+	tasks       *db.UploadTaskRepository
+	chunks      *db.UploadChunkRepository
+	maxUploadSz int64
+}
+
+func NewUploadTaskService(cfg *configs.UploadConfig, tasks *db.UploadTaskRepository, chunks *db.UploadChunkRepository) *UploadTaskService {
+	svc := &UploadTaskService{cfg: cfg, tasks: tasks, chunks: chunks}
+	if cfg != nil && cfg.MaxImageSizeMB > 0 {
+		svc.maxUploadSz = cfg.MaxImageSizeMB * 1024 * 1024
+	}
+	return svc
+}
+
+func (s *UploadTaskService) InitUpload(ctx context.Context, in InitInput) (*model.UploadTask, error) {
+	if s == nil {
+		return nil, errors.New("upload task service is required")
+	}
+	if in.UserID <= 0 {
+		return nil, errors.New("user id is required")
+	}
+	if in.FileName == "" {
+		return nil, errors.New("file name is required")
+	}
+	if in.FileSize <= 0 {
+		return nil, errors.New("file size is required")
+	}
+	if s.maxUploadSz > 0 && in.FileSize > s.maxUploadSz {
+		return nil, errors.New("file too large")
+	}
+	if s.tasks == nil {
+		return nil, errors.New("upload task repository is required")
+	}
+
+	chunkSize := in.ChunkSize
+	if chunkSize <= 0 {
+		chunkSize = defaultChunkSizeBytes
+	}
+	chunkCount := int(math.Ceil(float64(in.FileSize) / float64(chunkSize)))
+	if chunkCount <= 0 {
+		chunkCount = 1
+	}
+
+	task := &model.UploadTask{
+		UploadID:   uuid.NewString(),
+		UserID:     in.UserID,
+		BizType:    normalizeBizType(in.BizType),
+		BizID:      in.BizID,
+		FileName:   in.FileName,
+		FileSize:   in.FileSize,
+		ChunkSize:  chunkSize,
+		ChunkCount: chunkCount,
+		Status:     model.UploadTaskStatusUploading,
+		Sha256:     in.Sha256,
+		ExpiresAt:  time.Now().Add(24 * time.Hour).UTC(),
+	}
+	if err := s.tasks.Create(ctx, task); err != nil {
+		return nil, err
+	}
+	return task, nil
+}
+
+func (s *UploadTaskService) GetUpload(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, []model.UploadChunk, error) {
+	if s == nil || s.tasks == nil || s.chunks == nil {
+		return nil, nil, errors.New("upload task service is not ready")
+	}
+	task, err := s.tasks.GetByUploadID(ctx, uploadID, userID)
+	if err != nil {
+		return nil, nil, err
+	}
+	chunks, err := s.chunks.ListByUploadID(ctx, uploadID)
+	if err != nil {
+		return nil, nil, err
+	}
+	return task, chunks, nil
+}
+
+func (s *UploadTaskService) CancelUpload(ctx context.Context, uploadID string, userID int64) error {
+	return s.updateStatus(ctx, uploadID, userID, model.UploadTaskStatusCancelled)
+}
+
+func (s *UploadTaskService) CompleteUpload(ctx context.Context, uploadID string, userID int64) error {
+	return s.updateStatus(ctx, uploadID, userID, model.UploadTaskStatusCompleted)
+}
+
+func (s *UploadTaskService) updateStatus(ctx context.Context, uploadID string, userID int64, status string) error {
+	if s == nil || s.tasks == nil {
+		return errors.New("upload task repository is required")
+	}
+	task, err := s.tasks.GetByUploadID(ctx, uploadID, userID)
+	if err != nil {
+		return err
+	}
+	return s.tasks.UpdateProgress(ctx, task.UploadID, task.UserID, task.UploadedChunks, status)
+}
diff --git a/services/media-service/internal/service/upload_task_service_test.go b/services/media-service/internal/service/upload_task_service_test.go
new file mode 100644
index 0000000..5f18088
--- /dev/null
+++ b/services/media-service/internal/service/upload_task_service_test.go
@@ -0,0 +1,26 @@
+package service
+
+import (
+	"context"
+	"strings"
+	"testing"
+
+	"github.com/Loe1210/personal-site/configs"
+	"github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
+)
+
+func TestInitUploadRejectsTooLargeFile(t *testing.T) {
+	svc := NewUploadTaskService(&configs.UploadConfig{MaxImageSizeMB: 1}, db.NewUploadTaskRepository(nil), db.NewUploadChunkRepository(nil))
+
+	_, err := svc.InitUpload(context.Background(), InitInput{
+		UserID:   1,
+		FileName: "large.png",
+		FileSize: 2 * 1024 * 1024,
+	})
+	if err == nil {
+		t.Fatal("expected too large file to be rejected")
+	}
+	if !strings.Contains(err.Error(), "too large") {
+		t.Fatalf("expected too large error, got %v", err)
+	}
+}
```
