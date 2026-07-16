package model

import "time"

type FileRecord struct {
	ID           int64     `json:"id"`
	UploadID     string    `json:"upload_id"`
	OriginalName string    `json:"original_name"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Path         string    `json:"path"`
	ContentType  string    `json:"content_type"`
	Size         int64     `json:"size"`
	Sha256       string    `json:"sha256"`
	BizType      string    `json:"biz_type"`
	BizID        string    `json:"biz_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type UploadInput struct {
	FileName    string
	Content     []byte
	ContentType string
	BizType     string
}
