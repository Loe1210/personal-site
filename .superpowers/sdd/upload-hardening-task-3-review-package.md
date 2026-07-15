# Upload Hardening Task 3 Review Package

- Base: `efede0a`
- Head: `14ef305`
- Task: Chunk reception, tmp area, and resume

## Commits
14ef305 feat(media-service): add streamed chunk uploads

## Diff Stat
```text
 .superpowers/sdd/upload-hardening-task-3-brief.md  |  53 +++++++++
 services/media-service/biz/router.go               |   4 +-
 services/media-service/biz/upload/handler.go       |  46 +++++++-
 services/media-service/biz/upload/route.go         |   7 +-
 services/media-service/cmd/main.go                 |   3 +-
 services/media-service/cmd/router.go               |   4 +-
 .../internal/dal/storage/tmp_storage.go            |  86 +++++++++++++++
 .../internal/dal/storage/tmp_storage_test.go       |  34 ++++++
 .../internal/service/chunk_service.go              | 118 +++++++++++++++++++++
 .../internal/service/chunk_service_test.go         |  92 ++++++++++++++++
 .../internal/service/media_service.go              |   7 ++
 11 files changed, 443 insertions(+), 11 deletions(-)
```

## Full Diff
```diff
diff --git a/.superpowers/sdd/upload-hardening-task-3-brief.md b/.superpowers/sdd/upload-hardening-task-3-brief.md
new file mode 100644
index 0000000..279202d
--- /dev/null
+++ b/.superpowers/sdd/upload-hardening-task-3-brief.md
@@ -0,0 +1,53 @@
+# Task 3 Brief - Chunk reception, tmp area, and resume
+
+You are implementing Task 3 of the local file upload hardening plan in `docs/superpowers/plans/2026-07-15-local-file-upload-hardening.md`.
+
+## Where this task fits
+Task 1 added upload task/chunk persistence. Task 2 added the upload task state machine and route wiring. This task must add streamed chunk reception into a local temporary area and support resume bookkeeping.
+
+## Requirements
+- Use the repositories and models from Tasks 1 and 2.
+- Keep local disk storage only.
+- Do not read the whole large file into memory at once.
+- Store chunk bytes directly into the temporary upload area.
+- Record chunk metadata only after the write succeeds.
+- Keep the existing small-file upload flow intact.
+
+## Files in scope
+- Create: `services/media-service/internal/dal/storage/tmp_storage.go`
+- Create: `services/media-service/internal/dal/storage/tmp_storage_test.go`
+- Modify: `services/media-service/internal/service/media_service.go`
+- Create: `services/media-service/internal/service/chunk_service.go`
+- Create: `services/media-service/internal/service/chunk_service_test.go`
+- Modify: `services/media-service/biz/upload/handler.go`
+- Modify: `services/media-service/biz/upload/route.go`
+
+## Context you should assume
+- Task 1 and Task 2 are already committed and should not be reverted.
+- Upload task ownership comes from repository lookups using `upload_id` plus `user_id`.
+- The temp path must live under the local upload tmp directory, not in a shared object store.
+- You are not alone in the codebase. Do not revert or overwrite edits made by others.
+
+## Test-first target
+Start with the failing test from the plan:
+
+```go
+func TestChunkServiceWritesChunkToTmpPath(t *testing.T) {
+    // upload one chunk, assert temp file exists and metadata is recorded
+}
+```
+
+Then make it pass with the smallest change set that preserves the plan.
+
+## Report file
+Write your full report to:
+`.superpowers/sdd/upload-hardening-task-3-report.md`
+
+Report status should include:
+- what you implemented
+- tests run and results
+- any concerns or follow-up notes
+
+## Important
+- Commit your changes when done and include the commit hash in the report.
+- Keep the change scoped to Task 3 only.
diff --git a/services/media-service/biz/router.go b/services/media-service/biz/router.go
index eeed41f..7c4e7c3 100644
--- a/services/media-service/biz/router.go
+++ b/services/media-service/biz/router.go
@@ -1,13 +1,13 @@
 package biz
 
 import (
 	"github.com/cloudwego/hertz/pkg/app/server"
 
 	"github.com/Loe1210/personal-site/services/media-service/biz/upload"
 	"github.com/Loe1210/personal-site/services/media-service/internal/service"
 )
 
-func RegisterRoutes(hertz *server.Hertz, media *service.Service, uploadTasks *service.UploadTaskService) {
-	upload.RegisterRoutes(hertz, media)
+func RegisterRoutes(hertz *server.Hertz, media *service.Service, uploadTasks *service.UploadTaskService, chunks *service.ChunkService) {
+	upload.RegisterRoutes(hertz, media, chunks)
 	upload.RegisterTaskRoutes(hertz, uploadTasks)
 }
diff --git a/services/media-service/biz/upload/handler.go b/services/media-service/biz/upload/handler.go
index e16182f..b0ad55a 100644
--- a/services/media-service/biz/upload/handler.go
+++ b/services/media-service/biz/upload/handler.go
@@ -7,24 +7,25 @@ import (
 
 	"github.com/cloudwego/hertz/pkg/app"
 	"github.com/cloudwego/hertz/pkg/protocol/consts"
 
 	"github.com/Loe1210/personal-site/services/media-service/internal/model"
 	"github.com/Loe1210/personal-site/services/media-service/internal/service"
 )
 
 type Handler struct {
 	service *service.Service
+	chunks  *service.ChunkService
 }
 
-func NewHandler(service *service.Service) *Handler {
-	return &Handler{service: service}
+func NewHandler(service *service.Service, chunks *service.ChunkService) *Handler {
+	return &Handler{service: service, chunks: chunks}
 }
 
 func (h *Handler) Upload(ctx context.Context, c *app.RequestContext) {
 	header, err := c.FormFile("file")
 	if err != nil {
 		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20009, "message": "upload file is required"})
 		return
 	}
 	file, err := header.Open()
 	if err != nil {
@@ -43,23 +44,62 @@ func (h *Handler) Upload(ctx context.Context, c *app.RequestContext) {
 		ContentType: string(header.Header.Get("Content-Type")),
 		BizType:     c.PostForm("biz_type"),
 	})
 	if err != nil {
 		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20010, "message": err.Error()})
 		return
 	}
 	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": record})
 }
 
+func (h *Handler) UploadChunk(ctx context.Context, c *app.RequestContext) {
+	userID, err := parseUploadUserID(c)
+	if err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20030, "message": err.Error()})
+		return
+	}
+	chunkIndex, err := strconv.Atoi(c.Param("chunk_index"))
+	if err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20031, "message": "invalid chunk index"})
+		return
+	}
+	chunk, err := h.chunks.UploadChunk(ctx, service.ChunkInput{
+		UserID:     userID,
+		UploadID:   c.Param("upload_id"),
+		ChunkIndex: chunkIndex,
+		Body:       c.RequestBodyStream(),
+	})
+	if err != nil {
+		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20032, "message": err.Error()})
+		return
+	}
+	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": chunk})
+}
+
 func (h *Handler) GetFile(ctx context.Context, c *app.RequestContext) {
 	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
 	if err != nil {
 		c.JSON(consts.StatusBadRequest, map[string]any{"code": 20013, "message": "invalid file id"})
 		return
 	}
 	record, err := h.service.GetFile(ctx, id)
 	if err != nil {
 		c.JSON(consts.StatusNotFound, map[string]any{"code": 20014, "message": "file not found"})
 		return
 	}
 	c.JSON(consts.StatusOK, map[string]any{"code": 0, "message": "success", "data": record})
-}
\ No newline at end of file
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
diff --git a/services/media-service/biz/upload/route.go b/services/media-service/biz/upload/route.go
index 7d8b8be..1fe5b8e 100644
--- a/services/media-service/biz/upload/route.go
+++ b/services/media-service/biz/upload/route.go
@@ -1,13 +1,14 @@
 package upload
 
 import (
 	"github.com/cloudwego/hertz/pkg/app/server"
 
 	"github.com/Loe1210/personal-site/services/media-service/internal/service"
 )
 
-func RegisterRoutes(hertz *server.Hertz, media *service.Service) {
-	handler := NewHandler(media)
+func RegisterRoutes(hertz *server.Hertz, media *service.Service, chunks *service.ChunkService) {
+	handler := NewHandler(media, chunks)
 	hertz.POST("/upload", handler.Upload)
 	hertz.GET("/files/:id", handler.GetFile)
-}
\ No newline at end of file
+	hertz.POST("/upload/tasks/:upload_id/chunks/:chunk_index", handler.UploadChunk)
+}
diff --git a/services/media-service/cmd/main.go b/services/media-service/cmd/main.go
index bfa77ef..460c69b 100644
--- a/services/media-service/cmd/main.go
+++ b/services/media-service/cmd/main.go
@@ -31,16 +31,17 @@ func main() {
 	database, err := db.Open(cfg.MySQL)
 	if err != nil {
 		log.Fatal(err)
 	}
 	if err := db.Migrate(database); err != nil {
 		log.Fatal(err)
 	}
 	store := storage.NewLocalStorage(cfg.Upload.RootDir, cfg.Upload.PublicBasePath)
 	fileRepo := db.NewFileRepository(database)
 	uploadTasks := service.NewUploadTaskService(&cfg.Upload, db.NewUploadTaskRepository(database), db.NewUploadChunkRepository(database))
+	chunks := service.NewChunkService(db.NewUploadTaskRepository(database), db.NewUploadChunkRepository(database), storage.NewTmpStorage(""))
 	media := service.NewMediaService(store, fileRepo)
 	startMediaRPCServer(cfg.RPC.Port, kitexmediahandler.NewHandler(media))
-	h := newRouter(media, uploadTasks, configs.GetServerAddr())
+	h := newRouter(media, uploadTasks, chunks, configs.GetServerAddr())
 	log.Printf("media-service listening on %s", configs.GetServerAddr())
 	h.Spin()
 }
diff --git a/services/media-service/cmd/router.go b/services/media-service/cmd/router.go
index c7321a1..df7e0a3 100644
--- a/services/media-service/cmd/router.go
+++ b/services/media-service/cmd/router.go
@@ -1,14 +1,14 @@
 package main
 
 import (
 	"github.com/cloudwego/hertz/pkg/app/server"
 
 	"github.com/Loe1210/personal-site/services/media-service/biz"
 	"github.com/Loe1210/personal-site/services/media-service/internal/service"
 )
 
-func newRouter(media *service.Service, uploadTasks *service.UploadTaskService, addr string) *server.Hertz {
+func newRouter(media *service.Service, uploadTasks *service.UploadTaskService, chunks *service.ChunkService, addr string) *server.Hertz {
 	h := server.Default(server.WithHostPorts(addr))
-	biz.RegisterRoutes(h, media, uploadTasks)
+	biz.RegisterRoutes(h, media, uploadTasks, chunks)
 	return h
 }
diff --git a/services/media-service/internal/dal/storage/tmp_storage.go b/services/media-service/internal/dal/storage/tmp_storage.go
new file mode 100644
index 0000000..f4ff5ea
--- /dev/null
+++ b/services/media-service/internal/dal/storage/tmp_storage.go
@@ -0,0 +1,86 @@
+package storage
+
+import (
+	"crypto/sha256"
+	"encoding/hex"
+	"errors"
+	"fmt"
+	"io"
+	"os"
+	"path/filepath"
+	"strings"
+)
+
+type TmpStorage struct {
+	rootDir string
+}
+
+func NewTmpStorage(rootDir string) *TmpStorage {
+	if strings.TrimSpace(rootDir) == "" {
+		rootDir = "static/uploads/tmp"
+	}
+	return &TmpStorage{rootDir: rootDir}
+}
+
+func (s *TmpStorage) SaveChunk(uploadID string, chunkIndex int, content io.Reader) (string, int64, string, error) {
+	if s == nil {
+		return "", 0, "", errors.New("tmp storage is required")
+	}
+	if strings.TrimSpace(uploadID) == "" {
+		return "", 0, "", errors.New("upload id is required")
+	}
+	if chunkIndex < 0 {
+		return "", 0, "", errors.New("chunk index is required")
+	}
+	if content == nil {
+		return "", 0, "", errors.New("chunk content is required")
+	}
+
+	dir := filepath.Join(s.rootDir, uploadID)
+	if err := os.MkdirAll(dir, 0o755); err != nil {
+		return "", 0, "", err
+	}
+
+	storageName := fmt.Sprintf("chunk_%06d.part", chunkIndex)
+	tempPath := filepath.Join(dir, storageName+".tmp")
+	finalPath := filepath.Join(dir, storageName)
+	file, err := os.Create(tempPath)
+	if err != nil {
+		return "", 0, "", err
+	}
+
+	hash := sha256.New()
+	written, copyErr := io.Copy(io.MultiWriter(file, hash), content)
+	closeErr := file.Close()
+	if copyErr != nil {
+		_ = os.Remove(tempPath)
+		return "", 0, "", copyErr
+	}
+	if closeErr != nil {
+		_ = os.Remove(tempPath)
+		return "", 0, "", closeErr
+	}
+	if err := os.Rename(tempPath, finalPath); err != nil {
+		_ = os.Remove(tempPath)
+		return "", 0, "", err
+	}
+
+	return filepath.ToSlash(filepath.Join(uploadID, storageName)), written, hex.EncodeToString(hash.Sum(nil)), nil
+}
+
+func (s *TmpStorage) RemoveChunk(storagePath string) error {
+	if s == nil {
+		return errors.New("tmp storage is required")
+	}
+	if strings.TrimSpace(storagePath) == "" {
+		return nil
+	}
+	return os.Remove(filepath.Join(s.rootDir, filepath.FromSlash(storagePath)))
+}
+
+func (s *TmpStorage) Resolve(storagePath string) string {
+	if s == nil {
+		return ""
+	}
+	return filepath.Join(s.rootDir, filepath.FromSlash(storagePath))
+}
diff --git a/services/media-service/internal/dal/storage/tmp_storage_test.go b/services/media-service/internal/dal/storage/tmp_storage_test.go
new file mode 100644
index 0000000..9b8395f
--- /dev/null
+++ b/services/media-service/internal/dal/storage/tmp_storage_test.go
@@ -0,0 +1,34 @@
+package storage
+
+import (
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+)
+
+func TestTmpStorageWritesChunkToTmpPath(t *testing.T) {
+	tmpDir := t.TempDir()
+	store := NewTmpStorage(tmpDir)
+
+	storagePath, size, digest, err := store.SaveChunk("upload-1", 2, strings.NewReader("hello chunk"))
+	if err != nil {
+		t.Fatalf("save chunk: %v", err)
+	}
+	if storagePath != "upload-1/chunk_000002.part" {
+		t.Fatalf("unexpected storage path: %q", storagePath)
+	}
+	if size != int64(len("hello chunk")) {
+		t.Fatalf("unexpected size: %d", size)
+	}
+	if digest == "" {
+		t.Fatal("expected digest to be populated")
+	}
+	data, err := os.ReadFile(filepath.Join(tmpDir, filepath.FromSlash(storagePath)))
+	if err != nil {
+		t.Fatalf("read chunk: %v", err)
+	}
+	if string(data) != "hello chunk" {
+		t.Fatalf("unexpected chunk content: %q", string(data))
+	}
+}
diff --git a/services/media-service/internal/service/chunk_service.go b/services/media-service/internal/service/chunk_service.go
new file mode 100644
index 0000000..dd08387
--- /dev/null
+++ b/services/media-service/internal/service/chunk_service.go
@@ -0,0 +1,118 @@
+package service
+
+import (
+	"context"
+	"errors"
+	"fmt"
+	"io"
+	"sort"
+	"strconv"
+	"strings"
+
+	db "github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
+	"github.com/Loe1210/personal-site/services/media-service/internal/model"
+)
+
+type ChunkInput struct {
+	UserID     int64
+	UploadID   string
+	ChunkIndex int
+	Body       io.Reader
+}
+
+type ChunkService struct {
+	tasks   *db.UploadTaskRepository
+	chunks  *db.UploadChunkRepository
+	storage ChunkStorage
+}
+
+func NewChunkService(tasks *db.UploadTaskRepository, chunks *db.UploadChunkRepository, storage ChunkStorage) *ChunkService {
+	return &ChunkService{tasks: tasks, chunks: chunks, storage: storage}
+}
+
+func (s *ChunkService) UploadChunk(ctx context.Context, in ChunkInput) (*model.UploadChunk, error) {
+	if s == nil {
+		return nil, errors.New("chunk service is required")
+	}
+	if s.tasks == nil || s.chunks == nil || s.storage == nil {
+		return nil, errors.New("chunk service dependencies are required")
+	}
+	if in.UserID <= 0 {
+		return nil, errors.New("user id is required")
+	}
+	if strings.TrimSpace(in.UploadID) == "" {
+		return nil, errors.New("upload id is required")
+	}
+	if in.ChunkIndex < 0 {
+		return nil, errors.New("chunk index is required")
+	}
+	if in.Body == nil {
+		return nil, errors.New("chunk body is required")
+	}
+
+	task, err := s.tasks.GetByUploadID(ctx, in.UploadID, in.UserID)
+	if err != nil {
+		return nil, err
+	}
+	if task.Status != model.UploadTaskStatusUploading {
+		return nil, fmt.Errorf("upload task is not active: %s", task.Status)
+	}
+	if in.ChunkIndex >= task.ChunkCount {
+		return nil, fmt.Errorf("chunk index %d out of range", in.ChunkIndex)
+	}
+
+	storagePath, size, digest, err := s.storage.SaveChunk(in.UploadID, in.ChunkIndex, in.Body)
+	if err != nil {
+		return nil, err
+	}
+
+	chunk := &model.UploadChunk{
+		UploadID:    in.UploadID,
+		ChunkIndex:  in.ChunkIndex,
+		Size:        size,
+		Sha256:      digest,
+		StoragePath: storagePath,
+	}
+	if err := s.chunks.Save(ctx, chunk); err != nil {
+		_ = s.storage.RemoveChunk(storagePath)
+		return nil, err
+	}
+
+	uploadedChunks := mergeUploadedChunks(task.UploadedChunks, in.ChunkIndex)
+	if err := s.tasks.UpdateProgress(ctx, task.UploadID, task.UserID, uploadedChunks, task.Status); err != nil {
+		_ = s.storage.RemoveChunk(storagePath)
+		return nil, err
+	}
+
+	return chunk, nil
+}
+
+func mergeUploadedChunks(current string, chunkIndex int) string {
+	parts := strings.Split(current, ",")
+	seen := make(map[int]struct{}, len(parts)+1)
+	indices := make([]int, 0, len(parts)+1)
+	for _, part := range parts {
+		part = strings.TrimSpace(part)
+		if part == "" {
+			continue
+		}
+		idx, err := strconv.Atoi(part)
+		if err != nil {
+			continue
+		}
+		if _, ok := seen[idx]; ok {
+			continue
+		}
+		seen[idx] = struct{}{}
+		indices = append(indices, idx)
+	}
+	if _, ok := seen[chunkIndex]; !ok {
+		indices = append(indices, chunkIndex)
+	}
+	sort.Ints(indices)
+	items := make([]string, 0, len(indices))
+	for _, idx := range indices {
+		items = append(items, strconv.Itoa(idx))
+	}
+	return strings.Join(items, ",")
+}
diff --git a/services/media-service/internal/service/chunk_service_test.go b/services/media-service/internal/service/chunk_service_test.go
new file mode 100644
index 0000000..ec03b8f
--- /dev/null
+++ b/services/media-service/internal/service/chunk_service_test.go
@@ -0,0 +1,92 @@
+package service
+
+import (
+	"context"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+
+	db "github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
+	"github.com/Loe1210/personal-site/services/media-service/internal/dal/storage"
+	"github.com/Loe1210/personal-site/services/media-service/internal/model"
+	"github.com/glebarez/sqlite"
+	"gorm.io/gorm"
+)
+
+func TestChunkServiceWritesChunkToTmpPath(t *testing.T) {
+	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
+	if err != nil {
+		t.Fatalf("open test database: %v", err)
+	}
+	if err := db.Migrate(database); err != nil {
+		t.Fatalf("migrate test database: %v", err)
+	}
+
+	tmpDir := t.TempDir()
+	tmpStorage := storage.NewTmpStorage(tmpDir)
+	taskRepo := db.NewUploadTaskRepository(database)
+	chunkRepo := db.NewUploadChunkRepository(database)
+	svc := NewChunkService(taskRepo, chunkRepo, tmpStorage)
+
+	ctx := context.Background()
+	task := &model.UploadTask{
+		UploadID:   "upload-1",
+		UserID:     42,
+		BizType:    "article",
+		BizID:      "article-9",
+		FileName:   "video.mp4",
+		FileSize:   8 * 1024 * 1024,
+		ChunkSize:  4 * 1024 * 1024,
+		ChunkCount: 2,
+		Status:     model.UploadTaskStatusUploading,
+		ExpiresAt:  time.Now().Add(time.Hour).UTC(),
+	}
+	if err := taskRepo.Create(ctx, task); err != nil {
+		t.Fatalf("create upload task: %v", err)
+	}
+
+	chunk, err := svc.UploadChunk(ctx, ChunkInput{
+		UserID:     task.UserID,
+		UploadID:   task.UploadID,
+		ChunkIndex: 1,
+		Body:       strings.NewReader("hello chunk"),
+	})
+	if err != nil {
+		t.Fatalf("upload chunk: %v", err)
+	}
+	if chunk.StoragePath != "upload-1/chunk_000001.part" {
+		t.Fatalf("unexpected storage path: %q", chunk.StoragePath)
+	}
+
+	data, err := os.ReadFile(filepath.Join(tmpDir, filepath.FromSlash(chunk.StoragePath)))
+	if err != nil {
+		t.Fatalf("read chunk file: %v", err)
+	}
+	if string(data) != "hello chunk" {
+		t.Fatalf("unexpected chunk content: %q", string(data))
+	}
+
+	reloaded, err := taskRepo.GetByUploadID(ctx, task.UploadID, task.UserID)
+	if err != nil {
+		t.Fatalf("reload upload task: %v", err)
+	}
+	if reloaded.UploadedChunks != "1" {
+		t.Fatalf("expected uploaded chunks to be 1, got %q", reloaded.UploadedChunks)
+	}
+	if reloaded.Status != model.UploadTaskStatusUploading {
+		t.Fatalf("expected status uploading, got %q", reloaded.Status)
+	}
+
+	stored, err := chunkRepo.ListByUploadID(ctx, task.UploadID)
+	if err != nil {
+		t.Fatalf("list upload chunks: %v", err)
+	}
+	if len(stored) != 1 {
+		t.Fatalf("expected one stored chunk, got %d", len(stored))
+	}
+	if stored[0].ChunkIndex != 1 || stored[0].StoragePath != chunk.StoragePath {
+		t.Fatalf("unexpected stored chunk: %+v", stored[0])
+	}
+}
diff --git a/services/media-service/internal/service/media_service.go b/services/media-service/internal/service/media_service.go
index ee7712a..5e052ef 100644
--- a/services/media-service/internal/service/media_service.go
+++ b/services/media-service/internal/service/media_service.go
@@ -1,24 +1,30 @@
 package service
 
 import (
 	"context"
 	"errors"
+	"io"
 	"strings"
 
 	"github.com/Loe1210/personal-site/services/media-service/internal/model"
 )
 
 type Storage interface {
 	Save(name string, content []byte) (string, error)
 }
 
+type ChunkStorage interface {
+	SaveChunk(uploadID string, chunkIndex int, content io.Reader) (storagePath string, size int64, sha256 string, err error)
+	RemoveChunk(storagePath string) error
+}
+
 type Repository interface {
 	Save(ctx context.Context, record *model.FileRecord) error
 	GetByID(ctx context.Context, id int64) (*model.FileRecord, error)
 }
 
 type Service struct {
 	storage Storage
 	repo    Repository
 }
 
@@ -66,18 +72,19 @@ func (s *Service) GetFile(ctx context.Context, id int64) (*model.FileRecord, err
 	return s.repo.GetByID(ctx, id)
 }
 
 func normalizeBizType(input string) string {
 	bizType := strings.ToLower(strings.TrimSpace(input))
 	if bizType == "" {
 		return "common"
 	}
 	return bizType
 }
+
 func isAllowedImageContentType(contentType string) bool {
 	switch strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0])) {
 	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml":
 		return true
 	default:
 		return false
 	}
 }
```
