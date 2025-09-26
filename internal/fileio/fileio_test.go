package fileio

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/nalgeon/be"
)

func TestExists(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(path, []byte{1, 2, 3}, 0444)
		be.Err(t, err, nil)
		be.True(t, Exists(path))
	})
	t.Run("does not exist", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "file.txt")
		be.Equal(t, Exists(path), false)
	})
}

func TestCopyFiles(t *testing.T) {
	// create a temporary directory for testing
	srcDir, err := os.MkdirTemp("", "src")
	be.Err(t, err, nil)
	defer func() { _ = os.RemoveAll(srcDir) }()

	// create a source file
	srcFile := filepath.Join(srcDir, "source.txt")
	err = os.WriteFile(srcFile, []byte("test data"), 0644)
	be.Err(t, err, nil)

	// specify the destination directory
	dstDir, err := os.MkdirTemp("", "dst")
	be.Err(t, err, nil)
	defer func() { _ = os.RemoveAll(dstDir) }()

	t.Run("copy", func(t *testing.T) {
		// call the CopyFiles function
		const perm = fs.FileMode(0444)
		pattern := filepath.Join(srcDir, "*.txt")
		err = CopyFiles(pattern, dstDir, perm)
		be.Err(t, err, nil)

		// verify that the file was copied correctly
		dstFile := filepath.Join(dstDir, "source.txt")
		fileInfo, err := os.Stat(dstFile)
		be.Err(t, err, nil)
		be.Equal(t, fileInfo.Mode(), perm)

		// read the contents of the copied file
		data, err := os.ReadFile(dstFile)
		be.Err(t, err, nil)

		// verify the contents of the copied file
		expected := []byte("test data")
		be.Equal(t, data, expected)
	})

	t.Run("skip existing", func(t *testing.T) {
		// existing file in the destination dir
		path := filepath.Join(dstDir, "existing.txt")
		err := os.WriteFile(path, []byte("v1"), 0444)
		be.Err(t, err, nil)

		// same file in the source dir
		path = filepath.Join(srcDir, "existing.txt")
		err = os.WriteFile(path, []byte("v2"), 0444)
		be.Err(t, err, nil)

		// copy files
		pattern := filepath.Join(srcDir, "*.txt")
		err = CopyFiles(pattern, dstDir, 0444)
		be.Err(t, err, nil)

		// verify that the new file was copied correctly
		newFile := filepath.Join(dstDir, "source.txt")
		_, err = os.Stat(newFile)
		be.Err(t, err, nil)

		// verify that the existing file remained unchanged
		existFile := filepath.Join(dstDir, "existing.txt")
		data, err := os.ReadFile(existFile)
		be.Err(t, err, nil)
		expected := []byte("v1")
		be.Equal(t, data, expected)
	})
}

func TestReadJson(t *testing.T) {
	type Person struct{ Name string }

	t.Run("valid", func(t *testing.T) {
		got, err := ReadJson[Person](filepath.Join("testdata", "valid.json"))
		be.Err(t, err, nil)
		want := Person{"alice"}
		be.Equal(t, got, want)
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ReadJson[Person](filepath.Join("testdata", "invalid.json"))
		be.Err(t, err)
	})
	t.Run("does not exist", func(t *testing.T) {
		_, err := ReadJson[Person](filepath.Join("testdata", "missing.json"))
		be.Err(t, err)
	})
}

func TestWriteFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "files")
	be.Err(t, err, nil)
	defer func() { _ = os.RemoveAll(dir) }()

	t.Run("create nested dirs", func(t *testing.T) {
		path := filepath.Join(dir, "a/b/c/file.txt")
		err = WriteFile(path, "hello", 0444)
		be.Err(t, err, nil)
		got, err := os.ReadFile(path)
		be.Err(t, err, nil)
		want := []byte("hello")
		be.Equal(t, got, want)
	})

	t.Run("text", func(t *testing.T) {
		path := filepath.Join(dir, "data.txt")
		err = WriteFile(path, "hello", 0444)
		be.Err(t, err, nil)
		got, err := os.ReadFile(path)
		be.Err(t, err, nil)
		want := []byte("hello")
		be.Equal(t, got, want)
	})

	t.Run("data-octet-stream", func(t *testing.T) {
		path := filepath.Join(dir, "data-1.bin")
		err = WriteFile(path, "data:application/octet-stream;base64,MTIz", 0444)
		be.Err(t, err, nil)
		got, err := os.ReadFile(path)
		be.Err(t, err, nil)
		want := []byte("123")
		be.Equal(t, got, want)
	})

	t.Run("data-base64", func(t *testing.T) {
		path := filepath.Join(dir, "data-2.bin")
		err = WriteFile(path, "data:;base64,MTIz", 0444)
		be.Err(t, err, nil)
		got, err := os.ReadFile(path)
		be.Err(t, err, nil)
		want := []byte("123")
		be.Equal(t, got, want)
	})

	t.Run("data-text-plain", func(t *testing.T) {
		path := filepath.Join(dir, "data-3.bin")
		err = WriteFile(path, "data:text/plain;,123", 0444)
		be.Err(t, err, nil)
		got, err := os.ReadFile(path)
		be.Err(t, err, nil)
		want := []byte("123")
		be.Equal(t, got, want)
	})

	t.Run("perm", func(t *testing.T) {
		const perm = 0444
		path := filepath.Join(dir, "perm.txt")
		err = WriteFile(path, "hello", perm)
		be.Err(t, err, nil)
		fileInfo, err := os.Stat(path)
		be.Err(t, err, nil)
		be.Equal(t, fileInfo.Mode().Perm(), perm)
	})

	t.Run("missing data-url separator", func(t *testing.T) {
		path := filepath.Join(dir, "data.bin")
		err = WriteFile(path, "data:application/octet-stream:MTIz", 0444)
		be.Err(t, err)
	})

	t.Run("invalid binary value", func(t *testing.T) {
		path := filepath.Join(dir, "data.bin")
		err = WriteFile(path, "data:;base64,12345", 0444)
		be.Err(t, err)
	})
}

func TestJoinDir(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		filename string
		want     string
		wantErr  bool
	}{
		{
			name:     "regular join",
			dir:      "/home/user",
			filename: "docs/report.txt",
			want:     "/home/user/docs/report.txt",
			wantErr:  false,
		},
		{
			name:     "join with dot",
			dir:      "/home/user",
			filename: ".",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "join with absolute path",
			dir:      "/home/user",
			filename: "/etc/passwd",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "join with parent directory",
			dir:      "/home/user",
			filename: "../user2/docs/report.txt",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "empty directory",
			dir:      "",
			filename: "report.txt",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			dir:      "/home/user",
			filename: "",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "directory with trailing slash",
			dir:      "/home/user/",
			filename: "docs/report.txt",
			want:     "/home/user/docs/report.txt",
			wantErr:  false,
		},
		{
			name:     "filename with leading slash",
			dir:      "/home/user",
			filename: "/docs/report.txt",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "root directory",
			dir:      "/",
			filename: "report.txt",
			want:     "/report.txt",
			wantErr:  false,
		},
		{
			name:     "dot dot slash filename",
			dir:      "/home/user",
			filename: "..",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JoinDir(tt.dir, tt.filename)
			if tt.wantErr {
				be.Err(t, err)
			}
			be.Equal(t, got, tt.want)
		})
	}
}

func TestMkdirTemp(t *testing.T) {
	t.Run("default permissions", func(t *testing.T) {
		const perm = 0755
		dir, err := MkdirTemp(perm)
		be.Err(t, err, nil)
		defer func() { _ = os.Remove(dir) }()

		info, err := os.Stat(dir)
		be.Err(t, err, nil)
		be.Equal(t, info.Mode().Perm(), perm)
	})

	t.Run("non-default permissions", func(t *testing.T) {
		const perm = 0777
		dir, err := MkdirTemp(perm)
		be.Err(t, err, nil)
		defer func() { _ = os.Remove(dir) }()

		info, err := os.Stat(dir)
		be.Err(t, err, nil)
		be.Equal(t, info.Mode().Perm(), perm)
	})
}
