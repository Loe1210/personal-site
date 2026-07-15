package model

import "time"

const (
	UploadTaskStatusUploading = "uploading"
	UploadTaskStatusCompleted = "completed"
	UploadTaskStatusCancelled = "cancelled"
	UploadTaskStatusFailed    = "failed"
)

type UploadTask struct {
	UploadID       string
	UserID         int64
	BizType        string
	BizID          string
	FileName       string
	FileSize       int64
	ChunkSize      int64
	ChunkCount     int
	UploadedChunks string
	Status         string
	Sha256         string
	ExpiresAt      time.Time
	LastError      string
	Version        int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
