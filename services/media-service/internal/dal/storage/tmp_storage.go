package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type TmpStorage struct {
	rootDir string
}

func NewTmpStorage(rootDir string) *TmpStorage {
	if strings.TrimSpace(rootDir) == "" {
		rootDir = "static/uploads/tmp"
	}
	return &TmpStorage{rootDir: rootDir}
}

func (s *TmpStorage) SaveChunk(uploadID string, chunkIndex int, content io.Reader) (string, int64, string, error) {
	if s == nil {
		return "", 0, "", errors.New("tmp storage is required")
	}
	if strings.TrimSpace(uploadID) == "" {
		return "", 0, "", errors.New("upload id is required")
	}
	if chunkIndex < 0 {
		return "", 0, "", errors.New("chunk index is required")
	}
	if content == nil {
		return "", 0, "", errors.New("chunk content is required")
	}

	dir := filepath.Join(s.rootDir, uploadID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", 0, "", err
	}

	storageName := fmt.Sprintf("chunk_%06d.part", chunkIndex)
	tempPath := filepath.Join(dir, storageName+".tmp")
	finalPath := filepath.Join(dir, storageName)
	file, err := os.Create(tempPath)
	if err != nil {
		return "", 0, "", err
	}

	hash := sha256.New()
	written, copyErr := io.Copy(io.MultiWriter(file, hash), content)
	closeErr := file.Close()
	if copyErr != nil {
		_ = os.Remove(tempPath)
		return "", 0, "", copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tempPath)
		return "", 0, "", closeErr
	}
	if err := os.Remove(finalPath); err != nil && !os.IsNotExist(err) {
		_ = os.Remove(tempPath)
		return "", 0, "", err
	}
	if err := os.Rename(tempPath, finalPath); err != nil {
		_ = os.Remove(tempPath)
		return "", 0, "", err
	}

	return filepath.ToSlash(filepath.Join(uploadID, storageName)), written, hex.EncodeToString(hash.Sum(nil)), nil
}

func (s *TmpStorage) BackupChunk(storagePath string) (string, bool, error) {
	if s == nil {
		return "", false, errors.New("tmp storage is required")
	}
	if strings.TrimSpace(storagePath) == "" {
		return "", false, nil
	}
	sourcePath := s.Resolve(storagePath)
	source, err := os.Open(sourcePath)
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	defer source.Close()

	backup, err := os.CreateTemp(filepath.Dir(sourcePath), filepath.Base(sourcePath)+".bak-*")
	if err != nil {
		return "", false, err
	}
	backupPath := backup.Name()
	_, copyErr := io.Copy(backup, source)
	closeErr := backup.Close()
	if copyErr != nil {
		_ = os.Remove(backupPath)
		return "", false, copyErr
	}
	if closeErr != nil {
		_ = os.Remove(backupPath)
		return "", false, closeErr
	}
	return backupPath, true, nil
}

func (s *TmpStorage) RestoreChunk(storagePath string, backupPath string) error {
	if s == nil {
		return errors.New("tmp storage is required")
	}
	if strings.TrimSpace(storagePath) == "" || strings.TrimSpace(backupPath) == "" {
		return nil
	}
	targetPath := s.Resolve(storagePath)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}
	if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
		_ = os.Remove(backupPath)
		return err
	}
	if err := os.Rename(backupPath, targetPath); err != nil {
		_ = os.Remove(backupPath)
		return err
	}
	return nil
}

func (s *TmpStorage) DiscardChunkBackup(backupPath string) error {
	if s == nil {
		return errors.New("tmp storage is required")
	}
	if strings.TrimSpace(backupPath) == "" {
		return nil
	}
	if err := os.Remove(backupPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *TmpStorage) RemoveChunk(storagePath string) error {
	if s == nil {
		return errors.New("tmp storage is required")
	}
	if strings.TrimSpace(storagePath) == "" {
		return nil
	}
	if err := os.Remove(filepath.Join(s.rootDir, filepath.FromSlash(storagePath))); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *TmpStorage) Resolve(storagePath string) string {
	if s == nil {
		return ""
	}
	return filepath.Join(s.rootDir, filepath.FromSlash(storagePath))
}

func (s *TmpStorage) RemoveUpload(uploadID string) error {
	if s == nil {
		return errors.New("tmp storage is required")
	}
	if strings.TrimSpace(uploadID) == "" {
		return nil
	}
	if err := os.RemoveAll(filepath.Join(s.rootDir, uploadID)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
