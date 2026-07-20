package upload

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	localstorage "github.com/Loe1210/personal-site/services/media-service/internal/dal/storage"
	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"github.com/Loe1210/personal-site/services/media-service/internal/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

func TestUploadChunkReadsBufferedRequestBody(t *testing.T) {
	body := []byte("chunk body that must not become empty")
	sum := sha256.Sum256(body)
	expectedSHA := hex.EncodeToString(sum[:])

	tasks := &chunkTaskStore{task: &model.UploadTask{
		UploadID:   "upload-1",
		UserID:     1,
		ChunkCount: 1,
		Status:     model.UploadTaskStatusUploading,
		Version:    1,
	}}
	chunks := &chunkRecordStore{}
	storage := localstorage.NewTmpStorage(t.TempDir())
	chunkSvc := service.NewChunkService(tasks, chunks, storage)

	router := route.NewEngine(config.NewOptions(nil))
	handler := NewHandler(nil, chunkSvc)
	router.POST("/upload/tasks/:upload_id/chunks/:chunk_index", func(ctx context.Context, c *app.RequestContext) {
		handler.UploadChunk(ctx, c)
	})

	recorder := ut.PerformRequest(router, "POST", "/upload/tasks/upload-1/chunks/0?user_id=1", &ut.Body{Body: bytes.NewReader(body), Len: len(body)}, ut.Header{Key: "Content-Type", Value: "application/octet-stream"})
	resp := recorder.Result()
	if resp.StatusCode() != consts.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	var payload struct {
		Code int `json:"code"`
		Data struct {
			Size   int64  `json:"Size"`
			Sha256 string `json:"Sha256"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Code != 0 {
		t.Fatalf("expected code 0, got %d: %s", payload.Code, string(resp.Body()))
	}
	if payload.Data.Size != int64(len(body)) {
		t.Fatalf("expected chunk size %d, got %d", len(body), payload.Data.Size)
	}
	if payload.Data.Sha256 != expectedSHA {
		t.Fatalf("expected chunk sha %s, got %s", expectedSHA, payload.Data.Sha256)
	}
}

type chunkTaskStore struct {
	task           *model.UploadTask
	uploadedChunks string
}

func (s *chunkTaskStore) GetByUploadID(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, error) {
	task := *s.task
	task.UploadID = uploadID
	task.UserID = userID
	task.UploadedChunks = s.uploadedChunks
	return &task, nil
}

func (s *chunkTaskStore) UpdateProgressGuarded(ctx context.Context, uploadID string, userID int64, uploadedChunks string, status string, expectedStatus string, expectedVersion int64) error {
	s.uploadedChunks = uploadedChunks
	return nil
}

type chunkRecordStore struct {
	records []model.UploadChunk
}

func (s *chunkRecordStore) Save(ctx context.Context, chunk *model.UploadChunk) error {
	s.records = append(s.records, *chunk)
	return nil
}

func (s *chunkRecordStore) Delete(ctx context.Context, uploadID string, chunkIndex int) error {
	for i, chunk := range s.records {
		if chunk.UploadID == uploadID && chunk.ChunkIndex == chunkIndex {
			s.records = append(s.records[:i], s.records[i+1:]...)
			return nil
		}
	}
	return nil
}

func (s *chunkRecordStore) ListByUploadID(ctx context.Context, uploadID string) ([]model.UploadChunk, error) {
	out := make([]model.UploadChunk, 0, len(s.records))
	for _, chunk := range s.records {
		if chunk.UploadID == uploadID {
			out = append(out, chunk)
		}
	}
	return out, nil
}
