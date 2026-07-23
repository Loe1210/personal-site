package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var files []string
	if err := filepath.WalkDir("idl", func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".proto" {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		fail(err.Error())
	}
	if len(files) == 0 {
		fail("no proto files found under idl")
	}
	for _, file := range files {
		contentBytes, err := os.ReadFile(file)
		if err != nil {
			fail(fmt.Sprintf("read %s: %v", file, err))
		}
		content := string(contentBytes)
		require(file, content, "syntax = \"proto3\";")
		require(file, content, "package ")
		require(file, content, "option go_package")
		if !strings.HasPrefix(filepath.ToSlash(file), "idl/base/") {
			require(file, content, "service ")
		}
	}
}

func require(file string, content string, token string) {
	if !strings.Contains(content, token) {
		fail(fmt.Sprintf("%s missing %q", file, token))
	}
}

func fail(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
