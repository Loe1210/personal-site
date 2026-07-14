package db

import (
	"context"
	"time"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"gorm.io/gorm"
)

type FileRecord struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	OriginalName string    `gorm:"column:original_name;type:varchar(255);not null"`
	URL          string    `gorm:"column:url;type:varchar(255);not null"`
	Path         string    `gorm:"column:path;type:varchar(512);not null"`
	ContentType  string    `gorm:"column:content_type;type:varchar(128);not null"`
	Size         int64     `gorm:"column:size;not null"`
	BizType      string    `gorm:"column:biz_type;type:varchar(64);not null;default:'common'"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (FileRecord) TableName() string {
	return "files"
}

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Save(ctx context.Context, record *model.FileRecord) error {
	modelRecord := &FileRecord{
		OriginalName: record.OriginalName,
		URL:          record.URL,
		Path:         record.Path,
		ContentType:  record.ContentType,
		Size:         record.Size,
		BizType:      record.BizType,
	}
	if err := r.db.WithContext(ctx).Create(modelRecord).Error; err != nil {
		return err
	}
	record.ID = modelRecord.ID
	record.CreatedAt = modelRecord.CreatedAt
	return nil
}

func (r *FileRepository) GetByID(ctx context.Context, id int64) (*model.FileRecord, error) {
	var modelRecord FileRecord
	if err := r.db.WithContext(ctx).First(&modelRecord, id).Error; err != nil {
		return nil, err
	}
	return &model.FileRecord{
		ID:           modelRecord.ID,
		OriginalName: modelRecord.OriginalName,
		URL:          modelRecord.URL,
		Path:         modelRecord.Path,
		ContentType:  modelRecord.ContentType,
		Size:         modelRecord.Size,
		BizType:      modelRecord.BizType,
		CreatedAt:    modelRecord.CreatedAt,
	}, nil
}