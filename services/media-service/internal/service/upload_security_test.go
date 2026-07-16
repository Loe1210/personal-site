package service

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestUploadSecurityRejectsWrongMagicBytes(t *testing.T) {
	err := ValidateUploadContent("image/png", []byte("not really a png"), "")
	if err == nil {
		t.Fatal("expected wrong image magic bytes to be rejected")
	}
}

func TestUploadSecurityAcceptsMatchingSHA256(t *testing.T) {
	body := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}
	sum := sha256.Sum256(body)
	if err := ValidateUploadContent("image/png", body, hex.EncodeToString(sum[:])); err != nil {
		t.Fatalf("expected valid png upload to pass: %v", err)
	}
}

func TestUploadSecurityRejectsHashMismatch(t *testing.T) {
	body := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10}
	if err := ValidateUploadContent("image/jpeg", body, "deadbeef"); err == nil {
		t.Fatal("expected hash mismatch to be rejected")
	}
}
