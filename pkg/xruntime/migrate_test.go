package xruntime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrationFilesReturnsSortedSQLFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "002_second.sql", "select 2;")
	writeFile(t, dir, "001_first.sql", "select 1;")
	writeFile(t, dir, "notes.md", "skip")

	files, err := MigrationFiles(dir)
	if err != nil {
		t.Fatalf("MigrationFiles returned error: %v", err)
	}
	if len(files) != 2 || files[0] == files[1] || !endsWith(files[0], "001_first.sql") {
		t.Fatalf("unexpected files: %#v", files)
	}
}

func writeFile(t *testing.T, dir string, name string, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func endsWith(value string, suffix string) bool {
	return strings.HasSuffix(filepath.ToSlash(value), suffix)
}
