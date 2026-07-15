package service

import (
	"context"
	"strings"
	"testing"

	"github.com/Loe1210/personal-site/configs"
	"github.com/Loe1210/personal-site/services/media-service/internal/dal/db"
)

func TestInitUploadRejectsTooLargeFile(t *testing.T) {
	svc := NewUploadTaskService(&configs.UploadConfig{MaxImageSizeMB: 1}, db.NewUploadTaskRepository(nil), db.NewUploadChunkRepository(nil))

	_, err := svc.InitUpload(context.Background(), InitInput{
		UserID:   1,
		FileName: "large.png",
		FileSize: 2 * 1024 * 1024,
	})
	if err == nil {
		t.Fatal("expected too large file to be rejected")
	}
	if !strings.Contains(err.Error(), "too large") {
		t.Fatalf("expected too large error, got %v", err)
	}
}
