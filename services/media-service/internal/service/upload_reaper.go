package service

import (
	"context"
	"time"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
)

type ExpiredUploadTaskStore interface {
	ListExpired(ctx context.Context, now time.Time, limit int) ([]model.UploadTask, error)
	Delete(ctx context.Context, uploadID string) error
}

type UploadTmpCleaner interface {
	RemoveUpload(uploadID string) error
}

type UploadReaper struct {
	tasks ExpiredUploadTaskStore
	tmp   UploadTmpCleaner
	limit int
}

func NewUploadReaper(tasks ExpiredUploadTaskStore, tmp UploadTmpCleaner, limit int) *UploadReaper {
	if limit <= 0 {
		limit = 100
	}
	return &UploadReaper{tasks: tasks, tmp: tmp, limit: limit}
}

func (r *UploadReaper) RunOnce(ctx context.Context, now time.Time) (int, error) {
	if r == nil || r.tasks == nil || r.tmp == nil {
		return 0, nil
	}
	tasks, err := r.tasks.ListExpired(ctx, now, r.limit)
	if err != nil {
		return 0, err
	}
	deleted := 0
	for _, task := range tasks {
		if err := ctx.Err(); err != nil {
			return deleted, err
		}
		if err := r.tmp.RemoveUpload(task.UploadID); err != nil {
			return deleted, err
		}
		if err := r.tasks.Delete(ctx, task.UploadID); err != nil {
			return deleted, err
		}
		deleted++
	}
	return deleted, nil
}

func (r *UploadReaper) Start(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 10 * time.Minute
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-ticker.C:
				_, _ = r.RunOnce(ctx, now.UTC())
			}
		}
	}()
}
