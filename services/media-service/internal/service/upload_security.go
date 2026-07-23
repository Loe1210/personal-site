package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

func ValidateUploadContent(contentType string, content []byte, expectedSHA256 string) error {
	mediaType := normalizeContentType(contentType)
	if !isAllowedImageContentType(mediaType) {
		return errors.New("only image uploads are allowed")
	}
	if len(content) == 0 {
		return errors.New("file content is required")
	}
	if !matchesDeclaredImageType(mediaType, content) {
		return errors.New("file content does not match declared image type")
	}
	if expected := strings.TrimSpace(expectedSHA256); expected != "" {
		sum := sha256.Sum256(content)
		if !strings.EqualFold(expected, hex.EncodeToString(sum[:])) {
			return errors.New("file hash mismatch")
		}
	}
	return nil
}

func normalizeContentType(contentType string) string {
	return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}

func matchesDeclaredImageType(contentType string, content []byte) bool {
	sniffed := http.DetectContentType(content)
	switch contentType {
	case "image/jpeg":
		return sniffed == "image/jpeg"
	case "image/png":
		return sniffed == "image/png"
	case "image/gif":
		return sniffed == "image/gif"
	case "image/webp":
		return len(content) >= 12 && bytes.Equal(content[0:4], []byte("RIFF")) && bytes.Equal(content[8:12], []byte("WEBP"))
	case "image/svg+xml":
		trimmed := strings.TrimSpace(string(content))
		return strings.HasPrefix(trimmed, "<svg") || strings.Contains(trimmed, "<svg")
	default:
		return false
	}
}

func ValidateUploadFileContent(contentType string, filePath string) error {
	mediaType := normalizeContentType(contentType)
	if !isAllowedImageContentType(mediaType) {
		return errors.New("only image uploads are allowed")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	limited, err := io.ReadAll(io.LimitReader(file, 4096))
	if err != nil {
		return err
	}
	if len(limited) == 0 {
		return errors.New("file content is required")
	}
	if !matchesDeclaredImageType(mediaType, limited) {
		return errors.New("file content does not match declared image type")
	}
	return nil
}
