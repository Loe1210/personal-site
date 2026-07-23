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
	if err := validateUploadID(uploadID); err != nil {
		return "", 0, "", err
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
	file, err := os.CreateTemp(dir, storageName+".tmp-*")
	if err != nil {
		return "", 0, "", err
	}
	tempPath := file.Name()
	finalPath := filepath.Join(dir, storageName)

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

	digest := hex.EncodeToString(hash.Sum(nil))
	if existingDigest, existingSize, exists, err := checksumFile(finalPath); err != nil {
		_ = os.Remove(tempPath)
		return "", 0, "", err
	} else if exists && existingSize == written && existingDigest == digest {
		_ = os.Remove(tempPath)
		return filepath.ToSlash(filepath.Join(uploadID, storageName)), written, digest, nil
	}

	if err := os.Remove(finalPath); err != nil && !os.IsNotExist(err) {
		_ = os.Remove(tempPath)
		return "", 0, "", err
	}
	if err := os.Rename(tempPath, finalPath); err != nil {
		_ = os.Remove(tempPath)
		return "", 0, "", err
	}

	return filepath.ToSlash(filepath.Join(uploadID, storageName)), written, digest, nil
}

func (s *TmpStorage) BackupChunk(storagePath string) (string, bool, error) {
	if s == nil {
		return "", false, errors.New("tmp storage is required")
	}
	if strings.TrimSpace(storagePath) == "" {
		return "", false, nil
	}
	sourcePath, err := s.resolveStoragePath(storagePath)
	if err != nil {
		return "", false, err
	}
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
	targetPath, err := s.resolveStoragePath(storagePath)
	if err != nil {
		_ = os.Remove(backupPath)
		return err
	}
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
	resolved, err := s.resolveStoragePath(storagePath)
	if err != nil {
		return err
	}
	if err := os.Remove(resolved); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *TmpStorage) Resolve(storagePath string) string {
	if s == nil {
		return ""
	}
	resolved, err := s.resolveStoragePath(storagePath)
	if err != nil {
		return ""
	}
	return resolved
}

func (s *TmpStorage) RemoveUpload(uploadID string) error {
	if s == nil {
		return errors.New("tmp storage is required")
	}
	if strings.TrimSpace(uploadID) == "" {
		return nil
	}
	if err := validateUploadID(uploadID); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(s.rootDir, uploadID)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *TmpStorage) resolveStoragePath(storagePath string) (string, error) {
	cleaned := filepath.Clean(filepath.FromSlash(storagePath))
	if cleaned == "." || filepath.IsAbs(cleaned) || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) || cleaned == ".." || filepath.VolumeName(cleaned) != "" {
		return "", errors.New("invalid storage path")
	}
	return filepath.Join(s.rootDir, cleaned), nil
}

func validateUploadID(uploadID string) error {
	trimmed := strings.TrimSpace(uploadID)
	if trimmed == "" {
		return errors.New("upload id is required")
	}
	if trimmed != uploadID || trimmed == "." || trimmed == ".." || strings.ContainsAny(trimmed, `/\:`) || filepath.Clean(trimmed) != trimmed {
		return errors.New("invalid upload id")
	}
	return nil
}

func checksumFile(path string) (string, int64, bool, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return "", 0, false, nil
	}
	if err != nil {
		return "", 0, false, err
	}
	defer file.Close()

	hash := sha256.New()
	size, err := io.Copy(hash, file)
	if err != nil {
		return "", 0, false, err
	}
	return hex.EncodeToString(hash.Sum(nil)), size, true, nil
}
