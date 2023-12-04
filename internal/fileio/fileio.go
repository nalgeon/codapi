// Package fileio provides high-level file operations.
package fileio

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies all files matching the pattern
// to the destination directory.
func CopyFiles(pattern string, dstDir string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		src, err := os.Open(match)
		if err != nil {
			return err
		}
		defer src.Close()

		dstFile := filepath.Join(dstDir, filepath.Base(match))
		dst, err := os.Create(dstFile)
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadJson reads the file and decodes it from JSON.
func ReadJson[T any](path string) (T, error) {
	var obj T
	data, err := os.ReadFile(path)
	if err != nil {
		return obj, err
	}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, err
}
