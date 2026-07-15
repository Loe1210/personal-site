package model

import "time"

type UploadChunk struct {
	UploadID    string
	ChunkIndex  int
	Size        int64
	Sha256      string
	StoragePath string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
