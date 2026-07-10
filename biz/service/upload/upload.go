package upload

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	dbmodel "github.com/Loe1210/personal-site/biz/dal/db"
	uploadmodel "github.com/Loe1210/personal-site/biz/model/upload"
	"github.com/Loe1210/personal-site/pkg/errno"
)

const timeLayout = "2006-01-02 15:04:05"

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format(timeLayout)
}

func toUploadModel(item *dbmodel.UploadFile) *uploadmodel.UploadFile {
	if item == nil {
		return nil
	}

	return &uploadmodel.UploadFile{
		FileID:    item.ID,
		FileName:  item.FileName,
		FileURL:   item.FileURL,
		FilePath:  item.FilePath,
		MimeType:  item.MimeType,
		Size:      item.Size,
		BizType:   item.BizType,
		CreatedAt: formatTime(item.CreatedAt),
	}
}

func GetUploadInfo(_ context.Context, req *uploadmodel.GetUploadInfoRequest) (*uploadmodel.GetUploadInfoResponse, error) {
	var record dbmodel.UploadFile
	if err := dbmodel.DB.First(&record, req.ID).Error; err != nil {
		return nil, nil
	}

	return &uploadmodel.GetUploadInfoResponse{
		Upload: toUploadModel(&record),
	}, nil
}

func validateImageHeader(header *multipart.FileHeader) error {
	if header == nil {
		return errno.BadRequest
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		return errno.BadRequest
	}

	if !strings.HasPrefix(contentType, "image/") {
		return errno.BadRequest
	}

	if header.Size <= 0 {
		return errno.BadRequest
	}

	const maxImageSize = 5 * 1024 * 1024
	if header.Size > maxImageSize {
		return errno.BadRequest
	}

	return nil
}

func UploadImage(_ context.Context, req *uploadmodel.UploadImageRequest, header *multipart.FileHeader) (*uploadmodel.UploadImageResponse, error) {
	if err := validateImageHeader(header); err != nil {
		return nil, err
	}

	bizType := strings.TrimSpace(req.BizType)
	if bizType == "" {
		bizType = "common"
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".bin"
	}

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	relativeDir := filepath.Join("static", "uploads", "images")
	relativePath := filepath.Join(relativeDir, fileName)
	fileURL := "/" + filepath.ToSlash(relativePath)

	if err := os.MkdirAll(relativeDir, 0o755); err != nil {
		return nil, errno.Internal
	}

	if err := saveUploadedFile(header, relativePath); err != nil {
		return nil, errno.Internal
	}

	record := &dbmodel.UploadFile{
		FileName: header.Filename,
		FileURL:  fileURL,
		FilePath: relativePath,
		MimeType: header.Header.Get("Content-Type"),
		Size:     header.Size,
		BizType:  bizType,
	}

	if err := dbmodel.DB.Create(record).Error; err != nil {
		return nil, errno.Internal
	}

	return &uploadmodel.UploadImageResponse{
		Upload: toUploadModel(record),
	}, nil
}

func saveUploadedFile(header *multipart.FileHeader, dst string) error {
	src, err := header.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
