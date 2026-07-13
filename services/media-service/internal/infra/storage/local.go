package storage

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type LocalStorage struct {
	rootDir        string
	publicBasePath string
}

func NewLocalStorage(rootDir string, publicBasePath string) *LocalStorage {
	if strings.TrimSpace(rootDir) == "" {
		rootDir = "static/uploads/images"
	}
	if strings.TrimSpace(publicBasePath) == "" {
		publicBasePath = "/static/uploads/images"
	}
	return &LocalStorage{
		rootDir:        rootDir,
		publicBasePath: "/" + strings.Trim(publicBasePath, "/"),
	}
}

func (s *LocalStorage) Save(name string, content []byte) (string, error) {
	ext := filepath.Ext(name)
	dateDir := time.Now().Format("20060102")
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dir := filepath.Join(s.rootDir, dateDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dir, fileName), content, 0o644); err != nil {
		return "", err
	}
	return path.Join(s.publicBasePath, dateDir, fileName), nil
}
