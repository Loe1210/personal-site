package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	uploadmodel "github.com/Loe1210/personal-site/biz/model/upload"
	"github.com/Loe1210/personal-site/configs"
	dbmodel "github.com/Loe1210/personal-site/dal/db"
	"github.com/Loe1210/personal-site/pkg/errno"
)

var (
	allowedUploadImageTypes = map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
		"image/gif":  ".gif",
	}
	bizTypeSanitizer = regexp.MustCompile(`[^a-z0-9_-]+`)
)

type validatedImageMeta struct {
	mimeType string
	ext      string
}

func uploadConfig() configs.UploadConfig {
	if configs.AppConfig == nil {
		_, _ = configs.Load("")
	}
	if configs.AppConfig == nil {
		return configs.UploadConfig{
			RootDir:        "static/uploads/images",
			PublicBasePath: "/static/uploads/images",
			MaxImageSizeMB: 5,
		}
	}
	return configs.AppConfig.Upload
}

func maxUploadImageSizeBytes() int64 {
	return uploadConfig().MaxImageSizeMB * 1024 * 1024
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

func normalizeBizType(input string) string {
	bizType := strings.ToLower(strings.TrimSpace(input))
	if bizType == "" {
		return "common"
	}

	bizType = bizTypeSanitizer.ReplaceAllString(bizType, "-")
	bizType = strings.Trim(bizType, "-")
	if bizType == "" {
		return "common"
	}
	return bizType
}

func validateAndDetectImage(header *multipart.FileHeader) (*validatedImageMeta, error) {
	if header == nil {
		return nil, errno.UploadFileMissing
	}
	if header.Size <= 0 {
		return nil, errno.UploadFileEmpty
	}
	if header.Size > maxUploadImageSizeBytes() {
		return nil, errno.UploadFileTooLarge
	}

	file, err := header.Open()
	if err != nil {
		return nil, errno.Internal
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, readErr := file.Read(buffer)
	if readErr != nil && readErr != io.EOF {
		return nil, errno.Internal
	}
	if n == 0 {
		return nil, errno.UploadFileEmpty
	}

	detectedType := http.DetectContentType(buffer[:n])
	ext, ok := allowedUploadImageTypes[detectedType]
	if !ok {
		return nil, errno.UploadFileTypeInvalid
	}

	return &validatedImageMeta{
		mimeType: detectedType,
		ext:      ext,
	}, nil
}

func UploadImage(_ context.Context, req *uploadmodel.UploadImageRequest, header *multipart.FileHeader) (*uploadmodel.UploadImageResponse, error) {
	meta, err := validateAndDetectImage(header)
	if err != nil {
		return nil, err
	}

	cfg := uploadConfig()
	bizType := normalizeBizType(req.BizType)
	dateDir := time.Now().Format("20060102")
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), meta.ext)
	relativeDir := filepath.Join(cfg.RootDir, bizType, dateDir)
	relativePath := filepath.Join(relativeDir, fileName)
	publicBasePath := "/" + strings.Trim(strings.TrimSpace(cfg.PublicBasePath), "/")
	fileURL := path.Join(publicBasePath, bizType, dateDir, fileName)

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
		MimeType: meta.mimeType,
		Size:     header.Size,
		BizType:  bizType,
	}

	if err := dbmodel.DB.Create(record).Error; err != nil {
		_ = os.Remove(relativePath)
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
