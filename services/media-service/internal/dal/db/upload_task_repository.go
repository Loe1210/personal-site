package db

import (
	"context"
	"errors"
	"time"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"gorm.io/gorm"
)

var ErrUploadTaskStateConflict = errors.New("upload task state changed")

type UploadTaskRecord struct {
	UploadID       string    `gorm:"column:upload_id;type:varchar(64);primaryKey"`
	UserID         int64     `gorm:"column:user_id;not null;index:idx_upload_tasks_user_status"`
	BizType        string    `gorm:"column:biz_type;type:varchar(64);not null;default:'common'"`
	BizID          string    `gorm:"column:biz_id;type:varchar(128);not null;default:''"`
	FileName       string    `gorm:"column:file_name;type:varchar(255);not null"`
	FileSize       int64     `gorm:"column:file_size;not null"`
	ChunkSize      int64     `gorm:"column:chunk_size;not null"`
	ChunkCount     int       `gorm:"column:chunk_count;not null"`
	UploadedChunks string    `gorm:"column:uploaded_chunks;type:varchar(4096);not null;default:''"`
	Status         string    `gorm:"column:status;type:varchar(32);not null;index:idx_upload_tasks_user_status"`
	Sha256         string    `gorm:"column:sha256;type:varchar(64);not null;default:''"`
	ExpiresAt      time.Time `gorm:"column:expires_at;not null;index:idx_upload_tasks_expires_at"`
	LastError      string    `gorm:"column:last_error;type:text"`
	Version        int64     `gorm:"column:version;not null;default:1"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (UploadTaskRecord) TableName() string {
	return "upload_tasks"
}

type UploadTaskRepository struct {
	db *gorm.DB
}

func NewUploadTaskRepository(db *gorm.DB) *UploadTaskRepository {
	return &UploadTaskRepository{db: db}
}

func (r *UploadTaskRepository) Create(ctx context.Context, task *model.UploadTask) error {
	record := uploadTaskToRecord(task)
	if record.Status == "" {
		record.Status = model.UploadTaskStatusUploading
	}
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return err
	}
	copyUploadTaskRecord(task, record)
	return nil
}

func (r *UploadTaskRepository) GetByUploadID(ctx context.Context, uploadID string, userID int64) (*model.UploadTask, error) {
	var record UploadTaskRecord
	if err := r.db.WithContext(ctx).
		Where("upload_id = ? AND user_id = ?", uploadID, userID).
		First(&record).Error; err != nil {
		return nil, err
	}
	return uploadTaskFromRecord(&record), nil
}

func (r *UploadTaskRepository) UpdateProgressGuarded(ctx context.Context, uploadID string, userID int64, uploadedChunks string, status string, expectedStatus string, expectedVersion int64) error {
	query := r.db.WithContext(ctx).Model(&UploadTaskRecord{}).
		Where("upload_id = ? AND user_id = ?", uploadID, userID)
	if expectedStatus != "" {
		query = query.Where("status = ?", expectedStatus)
	}
	if expectedVersion > 0 {
		query = query.Where("version = ?", expectedVersion)
	}
	result := query.Updates(map[string]any{
		"uploaded_chunks": uploadedChunks,
		"status":          status,
		"version":         gorm.Expr("version + 1"),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUploadTaskStateConflict
	}
	return nil
}

func uploadTaskToRecord(task *model.UploadTask) *UploadTaskRecord {
	return &UploadTaskRecord{
		UploadID:       task.UploadID,
		UserID:         task.UserID,
		BizType:        task.BizType,
		BizID:          task.BizID,
		FileName:       task.FileName,
		FileSize:       task.FileSize,
		ChunkSize:      task.ChunkSize,
		ChunkCount:     task.ChunkCount,
		UploadedChunks: task.UploadedChunks,
		Status:         task.Status,
		Sha256:         task.Sha256,
		ExpiresAt:      task.ExpiresAt,
		LastError:      task.LastError,
		Version:        task.Version,
		CreatedAt:      task.CreatedAt,
		UpdatedAt:      task.UpdatedAt,
	}
}

func uploadTaskFromRecord(record *UploadTaskRecord) *model.UploadTask {
	task := &model.UploadTask{}
	copyUploadTaskRecord(task, record)
	return task
}

func copyUploadTaskRecord(task *model.UploadTask, record *UploadTaskRecord) {
	task.UploadID = record.UploadID
	task.UserID = record.UserID
	task.BizType = record.BizType
	task.BizID = record.BizID
	task.FileName = record.FileName
	task.FileSize = record.FileSize
	task.ChunkSize = record.ChunkSize
	task.ChunkCount = record.ChunkCount
	task.UploadedChunks = record.UploadedChunks
	task.Status = record.Status
	task.Sha256 = record.Sha256
	task.ExpiresAt = record.ExpiresAt
	task.LastError = record.LastError
	task.Version = record.Version
	task.CreatedAt = record.CreatedAt
	task.UpdatedAt = record.UpdatedAt
}
