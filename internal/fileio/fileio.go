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

// Exists checks if the specified path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	// we need a double negation here, because
	// errors.Is(err, os.ErrExist)
	// does not work
	return !errors.Is(err, os.ErrNotExist)
}

// CopyFile copies all files matching the pattern
// to the destination directory. Does not overwrite existing file.
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
		if Exists(dstFile) {
			continue
		}

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
	if !strings.HasPrefix(content, "data:") {
		// text file
		data = []byte(content)
		return os.WriteFile(path, data, perm)
	}

	// data-url encoded file
	meta, encoded, found := strings.Cut(content, ",")
	if !found {
		return errors.New("invalid data-url encoding")
	}

	if !strings.HasSuffix(meta, "base64") {
		// no need to decode
		data = []byte(encoded)
		return os.WriteFile(path, data, perm)
	}

	// decode base64-encoded data
	data, err = base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

// JoinDir joins a directory path with a relative file path,
// making sure that the resulting path is still inside the directory.
// Returns an error otherwise.
func JoinDir(dir string, name string) (string, error) {
	if dir == "" {
		return "", errors.New("invalid dir")
	}

	cleanName := filepath.Clean(name)
	if cleanName == "" {
		return "", errors.New("invalid name")
	}
	if cleanName == "." || cleanName == "/" || filepath.IsAbs(cleanName) {
		return "", errors.New("invalid name")
	}

	path := filepath.Join(dir, cleanName)

	dirPrefix := filepath.Clean(dir)
	if dirPrefix != "/" {
		dirPrefix += string(os.PathSeparator)
	}
	if !strings.HasPrefix(path, dirPrefix) {
		return "", errors.New("invalid name")
	}

	return path, nil
}

// MkdirTemp creates a new temporary directory with given permissions
// and returns the pathname of the new directory.
func MkdirTemp(perm fs.FileMode) (string, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}
	err = os.Chmod(dir, perm)
	if err != nil {
		os.Remove(dir)
		return "", err
	}
	return dir, nil
}
