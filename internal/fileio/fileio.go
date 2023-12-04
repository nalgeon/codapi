// Package fileio provides high-level file operations.
package fileio

import (
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
