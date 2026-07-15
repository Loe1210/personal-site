package db

import (
	"context"
	"time"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
	"gorm.io/gorm"
)

type UploadChunkRecord struct {
	UploadID    string    `gorm:"column:upload_id;type:varchar(64);primaryKey"`
	ChunkIndex  int       `gorm:"column:chunk_index;primaryKey"`
	Size        int64     `gorm:"column:size;not null"`
	Sha256      string    `gorm:"column:sha256;type:varchar(64);not null;default:''"`
	StoragePath string    `gorm:"column:storage_path;type:varchar(512);not null"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (UploadChunkRecord) TableName() string {
	return "upload_chunks"
}

type UploadChunkRepository struct {
	db *gorm.DB
}

func NewUploadChunkRepository(db *gorm.DB) *UploadChunkRepository {
	return &UploadChunkRepository{db: db}
}

func (r *UploadChunkRepository) Save(ctx context.Context, chunk *model.UploadChunk) error {
	record := uploadChunkToRecord(chunk)
	if err := r.db.WithContext(ctx).Save(record).Error; err != nil {
		return err
	}
	copyUploadChunkRecord(chunk, record)
	return nil
}

func (r *UploadChunkRepository) ListByUploadID(ctx context.Context, uploadID string) ([]model.UploadChunk, error) {
	var records []UploadChunkRecord
	if err := r.db.WithContext(ctx).
		Where("upload_id = ?", uploadID).
		Order("chunk_index ASC").
		Find(&records).Error; err != nil {
		return nil, err
	}
	chunks := make([]model.UploadChunk, 0, len(records))
	for i := range records {
		chunks = append(chunks, *uploadChunkFromRecord(&records[i]))
	}
	return chunks, nil
}

func uploadChunkToRecord(chunk *model.UploadChunk) *UploadChunkRecord {
	return &UploadChunkRecord{
		UploadID:    chunk.UploadID,
		ChunkIndex:  chunk.ChunkIndex,
		Size:        chunk.Size,
		Sha256:      chunk.Sha256,
		StoragePath: chunk.StoragePath,
		CreatedAt:   chunk.CreatedAt,
		UpdatedAt:   chunk.UpdatedAt,
	}
}

func uploadChunkFromRecord(record *UploadChunkRecord) *model.UploadChunk {
	chunk := &model.UploadChunk{}
	copyUploadChunkRecord(chunk, record)
	return chunk
}

func copyUploadChunkRecord(chunk *model.UploadChunk, record *UploadChunkRecord) {
	chunk.UploadID = record.UploadID
	chunk.ChunkIndex = record.ChunkIndex
	chunk.Size = record.Size
	chunk.Sha256 = record.Sha256
	chunk.StoragePath = record.StoragePath
	chunk.CreatedAt = record.CreatedAt
	chunk.UpdatedAt = record.UpdatedAt
}
