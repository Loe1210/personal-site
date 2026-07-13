package domain

import "time"

type FileRecord struct {
	ID           int64
	OriginalName string
	URL          string
	Path         string
	ContentType  string
	Size         int64
	BizType      string
	CreatedAt    time.Time
}
