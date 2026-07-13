package xruntime

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func MigrationFiles(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, filepath.Join(root, entry.Name()))
	}
	sort.Strings(files)
	return files, nil
}
