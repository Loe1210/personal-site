package service

import (
	"context"
	"testing"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
)

type fakeStorage struct {
	url string
}

func (f *fakeStorage) Save(name string, content []byte) (string, error) {
	return f.url, nil
}

type fakeRepository struct {
	saved  *model.FileRecord
	record *model.FileRecord
}

func (f *fakeRepository) Save(ctx context.Context, record *model.FileRecord) error {
	copyRecord := *record
	f.saved = &copyRecord
	record.ID = 9
	return nil
}

func (f *fakeRepository) GetByID(ctx context.Context, id int64) (*model.FileRecord, error) {
	return f.record, nil
}

func TestUploadReturnsStoredRecord(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewMediaService(&fakeStorage{url: "/uploads/cover.png"}, repo)

	resp, err := svc.Upload(context.Background(), model.UploadInput{
		FileName:    "cover.png",
		Content:     []byte("png"),
		ContentType: "image/png",
	})
	if err != nil {
		t.Fatalf("Upload returned error: %v", err)
	}
	if resp.URL != "/uploads/cover.png" {
		t.Fatalf("expected upload URL to stay intact, got %q", resp.URL)
	}
	if resp.Path != "/uploads/cover.png" {
		t.Fatalf("expected upload path to match URL, got %q", resp.Path)
	}
	if resp.BizType != "common" {
		t.Fatalf("expected default biz type common, got %q", resp.BizType)
	}
	if repo.saved == nil {
		t.Fatal("expected repository save to be called")
	}
}

func TestGetFileDelegatesToRepository(t *testing.T) {
	expected := &model.FileRecord{ID: 42, URL: "/uploads/answer.png"}
	svc := NewMediaService(&fakeStorage{url: "/uploads/unused.png"}, &fakeRepository{record: expected})

	resp, err := svc.GetFile(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetFile returned error: %v", err)
	}
	if resp != expected {
		t.Fatalf("expected repository record to be returned")
	}
}
