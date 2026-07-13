package application

import (
	"context"
	"testing"
)

type fakeStorage struct{}

func (f *fakeStorage) Save(name string, content []byte) (string, error) {
	return "/uploads/" + name, nil
}

func TestUploadReturnsURL(t *testing.T) {
	svc := NewMediaService(&fakeStorage{}, nil)
	resp, err := svc.Upload(context.Background(), UploadInput{
		FileName: "cover.png",
		Content:  []byte("png"),
	})
	if err != nil {
		t.Fatalf("Upload returned error: %v", err)
	}
	if resp.URL == "" {
		t.Fatal("expected upload URL")
	}
}
