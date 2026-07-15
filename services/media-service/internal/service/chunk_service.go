package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	db "github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/media-service/internal/model"
)

type ChunkInput struct {
	UserID     int64
	UploadID   string
	ChunkIndex int
	Body       io.Reader
}

type ChunkService struct {
	tasks   *db.UploadTaskRepository
	chunks  *db.UploadChunkRepository
	storage ChunkStorage
}

func NewChunkService(tasks *db.UploadTaskRepository, chunks *db.UploadChunkRepository, storage ChunkStorage) *ChunkService {
	return &ChunkService{tasks: tasks, chunks: chunks, storage: storage}
}

func (s *ChunkService) UploadChunk(ctx context.Context, in ChunkInput) (*model.UploadChunk, error) {
	if s == nil {
		return nil, errors.New("chunk service is required")
	}
	if s.tasks == nil || s.chunks == nil || s.storage == nil {
		return nil, errors.New("chunk service dependencies are required")
	}
	if in.UserID <= 0 {
		return nil, errors.New("user id is required")
	}
	if strings.TrimSpace(in.UploadID) == "" {
		return nil, errors.New("upload id is required")
	}
	if in.ChunkIndex < 0 {
		return nil, errors.New("chunk index is required")
	}
	if in.Body == nil {
		return nil, errors.New("chunk body is required")
	}

	task, err := s.tasks.GetByUploadID(ctx, in.UploadID, in.UserID)
	if err != nil {
		return nil, err
	}
	if task.Status != model.UploadTaskStatusUploading {
		return nil, fmt.Errorf("upload task is not active: %s", task.Status)
	}
	if in.ChunkIndex >= task.ChunkCount {
		return nil, fmt.Errorf("chunk index %d out of range", in.ChunkIndex)
	}

	storagePath, size, digest, err := s.storage.SaveChunk(in.UploadID, in.ChunkIndex, in.Body)
	if err != nil {
		return nil, err
	}

	chunk := &model.UploadChunk{
		UploadID:    in.UploadID,
		ChunkIndex:  in.ChunkIndex,
		Size:        size,
		Sha256:      digest,
		StoragePath: storagePath,
	}
	if err := s.chunks.Save(ctx, chunk); err != nil {
		_ = s.storage.RemoveChunk(storagePath)
		return nil, err
	}

	uploadedChunks := mergeUploadedChunks(task.UploadedChunks, in.ChunkIndex)
	if err := s.tasks.UpdateProgress(ctx, task.UploadID, task.UserID, uploadedChunks, task.Status); err != nil {
		_ = s.storage.RemoveChunk(storagePath)
		return nil, err
	}

	return chunk, nil
}

func mergeUploadedChunks(current string, chunkIndex int) string {
	parts := strings.Split(current, ",")
	seen := make(map[int]struct{}, len(parts)+1)
	indices := make([]int, 0, len(parts)+1)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		if _, ok := seen[idx]; ok {
			continue
		}
		seen[idx] = struct{}{}
		indices = append(indices, idx)
	}
	if _, ok := seen[chunkIndex]; !ok {
		indices = append(indices, chunkIndex)
	}
	sort.Ints(indices)
	items := make([]string, 0, len(indices))
	for _, idx := range indices {
		items = append(items, strconv.Itoa(idx))
	}
	return strings.Join(items, ",")
}
