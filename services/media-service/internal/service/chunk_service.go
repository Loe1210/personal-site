package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
)

type ChunkInput struct {
	UserID     int64
	UploadID   string
	ChunkIndex int
	Body       io.Reader
}

type UploadTaskStore interface {
	GetByUploadID(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, error)
	UpdateProgressGuarded(ctx context.Context, uploadID string, userID int64, uploadedChunks string, status string, expectedStatus string, expectedVersion int64) error
}

type UploadChunkStore interface {
	Save(ctx context.Context, chunk *model.UploadChunk) error
	Delete(ctx context.Context, uploadID string, chunkIndex int) error
	ListByUploadID(ctx context.Context, uploadID string) ([]model.UploadChunk, error)
}

type ChunkService struct {
	tasks   UploadTaskStore
	chunks  UploadChunkStore
	storage ChunkStorage
}

func NewChunkService(tasks UploadTaskStore, chunks UploadChunkStore, storage ChunkStorage) *ChunkService {
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

	previousChunk, previousBackupPath, err := s.backupExistingChunk(ctx, in.UploadID, in.ChunkIndex)
	if err != nil {
		return nil, err
	}

	storagePath, size, digest, err := s.storage.SaveChunk(in.UploadID, in.ChunkIndex, in.Body)
	if err != nil {
		_ = s.restorePreviousChunk(ctx, in.UploadID, in.ChunkIndex, previousChunk, previousBackupPath, storagePath)
		return nil, err
	}

	chunk := &model.UploadChunk{
		UploadID:    in.UploadID,
		ChunkIndex:  in.ChunkIndex,
		Size:        size,
		Sha256:      digest,
		StoragePath: storagePath,
	}
	if previousChunk != nil {
		if err := s.chunks.Delete(ctx, in.UploadID, in.ChunkIndex); err != nil {
			_ = s.restorePreviousChunk(ctx, in.UploadID, in.ChunkIndex, previousChunk, previousBackupPath, storagePath)
			return nil, err
		}
	}
	if err := s.chunks.Save(ctx, chunk); err != nil {
		_ = s.restorePreviousChunk(ctx, in.UploadID, in.ChunkIndex, previousChunk, previousBackupPath, storagePath)
		return nil, err
	}

	task, err = s.tasks.GetByUploadID(ctx, in.UploadID, in.UserID)
	if err != nil {
		_ = s.restorePreviousChunk(ctx, in.UploadID, in.ChunkIndex, previousChunk, previousBackupPath, storagePath)
		return nil, err
	}
	if task.Status != model.UploadTaskStatusUploading {
		_ = s.restorePreviousChunk(ctx, in.UploadID, in.ChunkIndex, previousChunk, previousBackupPath, storagePath)
		return nil, fmt.Errorf("upload task is not active: %s", task.Status)
	}

	uploadedChunks := mergeUploadedChunks(task.UploadedChunks, in.ChunkIndex)
	if err := s.tasks.UpdateProgressGuarded(ctx, task.UploadID, task.UserID, uploadedChunks, task.Status, task.Status, task.Version); err != nil {
		rollbackErr := s.restorePreviousChunk(ctx, in.UploadID, in.ChunkIndex, previousChunk, previousBackupPath, storagePath)
		if rollbackErr != nil {
			return nil, fmt.Errorf("update progress failed: %w; rollback failed: %v", err, rollbackErr)
		}
		return nil, err
	}

	_ = s.storage.DiscardChunkBackup(previousBackupPath)
	return chunk, nil
}

func (s *ChunkService) backupExistingChunk(ctx context.Context, uploadID string, chunkIndex int) (*model.UploadChunk, string, error) {
	chunk, err := s.findChunk(ctx, uploadID, chunkIndex)
	if err != nil || chunk == nil {
		return chunk, "", err
	}
	backupPath, exists, err := s.storage.BackupChunk(chunk.StoragePath)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		return nil, "", nil
	}
	return chunk, backupPath, nil
}

func (s *ChunkService) findChunk(ctx context.Context, uploadID string, chunkIndex int) (*model.UploadChunk, error) {
	chunks, err := s.chunks.ListByUploadID(ctx, uploadID)
	if err != nil {
		return nil, err
	}
	for i := range chunks {
		if chunks[i].ChunkIndex == chunkIndex {
			return &chunks[i], nil
		}
	}
	return nil, nil
}

func (s *ChunkService) restorePreviousChunk(ctx context.Context, uploadID string, chunkIndex int, previousChunk *model.UploadChunk, previousBackupPath string, currentStoragePath string) error {
	if previousChunk == nil {
		rollbackErr := s.chunks.Delete(ctx, uploadID, chunkIndex)
		removeErr := s.storage.RemoveChunk(currentStoragePath)
		if rollbackErr != nil {
			return rollbackErr
		}
		return removeErr
	}
	if err := s.chunks.Delete(ctx, uploadID, chunkIndex); err != nil {
		_ = s.storage.RemoveChunk(currentStoragePath)
		_ = s.storage.DiscardChunkBackup(previousBackupPath)
		return err
	}
	if err := s.chunks.Save(ctx, previousChunk); err != nil {
		_ = s.storage.RemoveChunk(currentStoragePath)
		_ = s.storage.DiscardChunkBackup(previousBackupPath)
		return err
	}
	if err := s.storage.RestoreChunk(previousChunk.StoragePath, previousBackupPath); err != nil {
		return err
	}
	return nil
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
