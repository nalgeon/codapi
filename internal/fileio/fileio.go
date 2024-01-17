// Package fileio provides high-level file operations.
package fileio

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CopyFile copies all files matching the pattern
// to the destination directory.
func CopyFiles(pattern string, dstDir string, perm fs.FileMode) error {
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
		dst, err := os.OpenFile(dstFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
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

// WriteFile writes the file to disk.
// The content can be text or binary (encoded as a data URL),
// e.g. data:application/octet-stream;base64,MTIz
func WriteFile(path, content string, perm fs.FileMode) (err error) {
	var data []byte
	if strings.HasPrefix(content, "data:") {
		// data-url encoded file
		_, encoded, found := strings.Cut(content, ",")
		if !found {
			return errors.New("invalid data-url encoding")
		}
		data, err = base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return err
		}
	} else {
		// text file
		data = []byte(content)
	}
	return os.WriteFile(path, data, perm)
}
