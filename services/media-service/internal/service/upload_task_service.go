package service

import (
	"context"
	"errors"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Loe1210/personal-site/configs"
	db "github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"github.com/google/uuid"
)

const defaultChunkSizeBytes int64 = 4 * 1024 * 1024

type InitInput struct {
	UserID      int64
	FileName    string
	FileSize    int64
	ContentType string
	BizType     string
	BizID       string
	Sha256      string
	ChunkSize   int64
}

type UploadTaskRepository interface {
	Create(ctx context.Context, task *model.UploadTask) error
	GetByUploadID(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, error)
	UpdateProgressGuarded(ctx context.Context, uploadID string, userID int64, uploadedChunks string, status string, expectedStatus string, expectedVersion int64) error
}

type UploadCompletionStore interface {
	SaveRecordAndCompleteTask(ctx context.Context, task *model.UploadTask, record *model.FileRecord) error
}

type UploadTaskService struct {
	cfg            *configs.UploadConfig
	tasks          UploadTaskRepository
	chunks         *db.UploadChunkRepository
	maxUploadSz    int64
	merge          *MergeService
	completion     UploadCompletionStore
	imageProcessor *ImageProcessor
}

func NewUploadTaskService(cfg *configs.UploadConfig, tasks UploadTaskRepository, chunks *db.UploadChunkRepository) *UploadTaskService {
	svc := &UploadTaskService{cfg: cfg, tasks: tasks, chunks: chunks}
	if cfg != nil && cfg.MaxImageSizeMB > 0 {
		svc.maxUploadSz = cfg.MaxImageSizeMB * 1024 * 1024
	}
	return svc
}

func (s *UploadTaskService) InitUpload(ctx context.Context, in InitInput) (*model.UploadTask, error) {
	if s == nil {
		return nil, errors.New("upload task service is required")
	}
	if in.UserID <= 0 {
		return nil, errors.New("user id is required")
	}
	if in.FileName == "" {
		return nil, errors.New("file name is required")
	}
	if in.FileSize <= 0 {
		return nil, errors.New("file size is required")
	}
	if s.maxUploadSz > 0 && in.FileSize > s.maxUploadSz {
		return nil, errors.New("file too large")
	}
	if !isAllowedImageContentType(normalizeContentType(in.ContentType)) {
		return nil, errors.New("only image uploads are allowed")
	}
	if s.tasks == nil {
		return nil, errors.New("upload task repository is required")
	}

	chunkSize := in.ChunkSize
	if chunkSize <= 0 {
		chunkSize = defaultChunkSizeBytes
	}
	chunkCount := int(math.Ceil(float64(in.FileSize) / float64(chunkSize)))
	if chunkCount <= 0 {
		chunkCount = 1
	}

	task := &model.UploadTask{
		UploadID:    uuid.NewString(),
		UserID:      in.UserID,
		BizType:     normalizeBizType(in.BizType),
		BizID:       in.BizID,
		FileName:    in.FileName,
		FileSize:    in.FileSize,
		ContentType: in.ContentType,
		ChunkSize:   chunkSize,
		ChunkCount:  chunkCount,
		Status:      model.UploadTaskStatusUploading,
		Sha256:      in.Sha256,
		ExpiresAt:   time.Now().Add(24 * time.Hour).UTC(),
	}
	if err := s.tasks.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *UploadTaskService) GetUpload(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, []model.UploadChunk, error) {
	if s == nil || s.tasks == nil || s.chunks == nil {
		return nil, nil, errors.New("upload task service is not ready")
	}
	task, err := s.tasks.GetByUploadID(ctx, uploadID, userID)
	if err != nil {
		return nil, nil, err
	}
	chunks, err := s.chunks.ListByUploadID(ctx, uploadID)
	if err != nil {
		return nil, nil, err
	}
	return task, chunks, nil
}

func (s *UploadTaskService) CancelUpload(ctx context.Context, uploadID string, userID int64) error {
	return s.updateStatus(ctx, uploadID, userID, model.UploadTaskStatusCancelled)
}

func (s *UploadTaskService) ConfigureCompletion(merge *MergeService, completion UploadCompletionStore, processors ...*ImageProcessor) {
	s.merge = merge
	s.completion = completion
	if len(processors) > 0 {
		s.imageProcessor = processors[0]
		return
	}
	s.imageProcessor = NewImageProcessor()
}

func (s *UploadTaskService) CompleteUpload(ctx context.Context, uploadID string, userID int64) (*model.FileRecord, error) {
	if s == nil || s.tasks == nil || s.chunks == nil || s.merge == nil || s.completion == nil {
		return nil, errors.New("upload completion pipeline is not configured")
	}
	task, err := s.tasks.GetByUploadID(ctx, uploadID, userID)
	if err != nil {
		return nil, err
	}
	if task.Status != model.UploadTaskStatusUploading {
		return nil, errors.New("upload task is not active")
	}
	chunks, err := s.chunks.ListByUploadID(ctx, uploadID)
	if err != nil {
		return nil, err
	}
	result, err := s.merge.Merge(ctx, MergeInput{UploadID: task.UploadID, FileName: task.FileName, ExpectedSHA256: task.Sha256, ChunkCount: task.ChunkCount, Chunks: chunks})
	if err != nil {
		return nil, err
	}
	if result.Size != task.FileSize {
		s.cleanupCompletedFiles(result, "")
		return nil, errors.New("final file size mismatch")
	}
	if err := ValidateUploadFileContent(task.ContentType, result.FinalPath); err != nil {
		s.cleanupCompletedFiles(result, "")
		return nil, err
	}

	thumbnailURL, err := s.createThumbnail(result)
	if err != nil {
		s.cleanupCompletedFiles(result, "")
		return nil, err
	}
	record := &model.FileRecord{
		UploadID:     task.UploadID,
		OriginalName: task.FileName,
		URL:          result.PublicPath,
		ThumbnailURL: thumbnailURL,
		Path:         result.RelativePath,
		ContentType:  task.ContentType,
		Size:         result.Size,
		Sha256:       result.Sha256,
		BizType:      task.BizType,
		BizID:        task.BizID,
	}
	if err := s.completion.SaveRecordAndCompleteTask(ctx, task, record); err != nil {
		s.cleanupCompletedFiles(result, thumbnailURL)
		return nil, err
	}
	return record, nil
}

func (s *UploadTaskService) createThumbnail(result *MergeResult) (string, error) {
	if s.imageProcessor == nil || result == nil || result.FinalPath == "" {
		return "", nil
	}
	thumbRelative := thumbnailRelativePath(result.RelativePath)
	thumbPath := filepath.Join(filepath.Dir(result.FinalPath), "thumbs", filepath.Base(thumbRelative))
	created, err := s.imageProcessor.Process(result.FinalPath, thumbPath)
	if err != nil || !created {
		return "", err
	}
	return path.Join(path.Dir(result.PublicPath), "thumbs", filepath.Base(thumbRelative)), nil
}

func (s *UploadTaskService) cleanupCompletedFiles(result *MergeResult, thumbnailURL string) {
	if result == nil {
		return
	}
	if thumbnailURL != "" {
		thumbRelative := thumbnailRelativePath(result.RelativePath)
		thumbPath := filepath.Join(filepath.Dir(result.FinalPath), "thumbs", filepath.Base(thumbRelative))
		removeFileAndEmptyParents(filepath.Dir(result.FinalPath), thumbPath)
	}
	removeFileAndEmptyParents(filepath.Dir(filepath.Dir(result.FinalPath)), result.FinalPath)
}

func removeFileAndEmptyParents(root string, target string) {
	if target == "" {
		return
	}
	_ = os.Remove(target)
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return
	}
	for dir := filepath.Dir(target); dir != "." && dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		dirAbs, err := filepath.Abs(dir)
		if err != nil || dirAbs == rootAbs || !strings.HasPrefix(dirAbs, rootAbs+string(filepath.Separator)) {
			return
		}
		if err := os.Remove(dir); err != nil {
			return
		}
	}
}
func thumbnailRelativePath(relative string) string {
	name := path.Base(filepath.ToSlash(relative))
	ext := path.Ext(name)
	base := strings.TrimSuffix(name, ext)
	if base == "" {
		base = "thumbnail"
	}
	return path.Join(path.Dir(filepath.ToSlash(relative)), "thumbs", base+".jpg")
}

func (s *UploadTaskService) updateStatus(ctx context.Context, uploadID string, userID int64, status string) error {
	if s == nil || s.tasks == nil {
		return errors.New("upload task repository is required")
	}
	task, err := s.tasks.GetByUploadID(ctx, uploadID, userID)
	if err != nil {
		return err
	}
	return s.tasks.UpdateProgressGuarded(ctx, task.UploadID, task.UserID, task.UploadedChunks, status, task.Status, task.Version)
}
