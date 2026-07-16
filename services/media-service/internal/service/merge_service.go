package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Loe1210/personal-site/services/media-service/internal/model"
)

type MergeInput struct {
	UploadID       string
	FileName       string
	ExpectedSHA256 string
	ChunkCount     int
	Chunks         []model.UploadChunk
}

type MergeResult struct {
	RelativePath string
	PublicPath   string
	FinalPath    string
	Sha256       string
	Size         int64
}

type MergeService struct {
	tmpRoot    string
	finalRoot  string
	publicBase string
}

func NewMergeService(tmpRoot, finalRoot, publicBase string) *MergeService {
	return &MergeService{tmpRoot: tmpRoot, finalRoot: finalRoot, publicBase: "/" + strings.Trim(publicBase, "/")}
}

func (s *MergeService) Merge(ctx context.Context, in MergeInput) (*MergeResult, error) {
	if s == nil || in.UploadID == "" || in.FileName == "" || in.ChunkCount <= 0 || len(in.Chunks) != in.ChunkCount {
		return nil, errors.New("complete chunk metadata is required")
	}
	chunks := append([]model.UploadChunk(nil), in.Chunks...)
	sort.Slice(chunks, func(i, j int) bool { return chunks[i].ChunkIndex < chunks[j].ChunkIndex })

	dateDir := time.Now().Format("20060102")
	dir := filepath.Join(s.finalRoot, dateDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	tmp, err := os.CreateTemp(dir, "merge-*.tmp")
	if err != nil {
		return nil, err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	hash := sha256.New()
	var size int64
	for index, chunk := range chunks {
		if chunk.UploadID != in.UploadID || chunk.ChunkIndex != index || chunk.StoragePath == "" {
			_ = tmp.Close()
			return nil, errors.New("invalid chunk metadata")
		}
		if err := ctx.Err(); err != nil {
			_ = tmp.Close()
			return nil, err
		}
		source, err := os.Open(filepath.Join(s.tmpRoot, filepath.FromSlash(chunk.StoragePath)))
		if err != nil {
			_ = tmp.Close()
			return nil, err
		}
		written, copyErr := io.Copy(io.MultiWriter(tmp, hash), source)
		closeErr := source.Close()
		if copyErr != nil {
			_ = tmp.Close()
			return nil, copyErr
		}
		if closeErr != nil {
			_ = tmp.Close()
			return nil, closeErr
		}
		if written != chunk.Size {
			_ = tmp.Close()
			return nil, errors.New("chunk size mismatch")
		}
		size += written
	}

	actualSHA256 := hex.EncodeToString(hash.Sum(nil))
	if expected := strings.TrimSpace(in.ExpectedSHA256); expected != "" && !strings.EqualFold(expected, actualSHA256) {
		_ = tmp.Close()
		return nil, errors.New("final file hash mismatch")
	}
	if err := tmp.Close(); err != nil {
		return nil, err
	}

	baseName := strings.TrimSuffix(filepath.Base(tmpName), ".tmp") + strings.ToLower(filepath.Ext(in.FileName))
	relative := filepath.ToSlash(filepath.Join(dateDir, baseName))
	finalPath := filepath.Join(s.finalRoot, filepath.FromSlash(relative))
	if err := os.Rename(tmpName, finalPath); err != nil {
		return nil, err
	}

	return &MergeResult{RelativePath: relative, PublicPath: path.Join(s.publicBase, relative), FinalPath: finalPath, Size: size, Sha256: actualSHA256}, nil
}
