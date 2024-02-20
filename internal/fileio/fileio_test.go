package fileio

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExists(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "file.txt")
		err := os.WriteFile(path, []byte{1, 2, 3}, 0444)
		if err != nil {
			t.Fatal(err)
		}
		if !Exists(path) {
			t.Fatalf("Exists: %s does not exist", filepath.Base(path))
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "file.txt")
		if Exists(path) {
			t.Fatalf("Exists: %s should not exist", filepath.Base(path))
		}
	})
}

func TestCopyFiles(t *testing.T) {
	// create a temporary directory for testing
	srcDir, err := os.MkdirTemp("", "src")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	// create a source file
	srcFile := filepath.Join(srcDir, "source.txt")
	err = os.WriteFile(srcFile, []byte("test data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// specify the destination directory
	dstDir, err := os.MkdirTemp("", "dst")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dstDir)

	t.Run("copy", func(t *testing.T) {
		// call the CopyFiles function
		const perm = fs.FileMode(0444)
		pattern := filepath.Join(srcDir, "*.txt")
		err = CopyFiles(pattern, dstDir, perm)
		if err != nil {
			t.Fatal(err)
		}

		// verify that the file was copied correctly
		dstFile := filepath.Join(dstDir, "source.txt")
		fileInfo, err := os.Stat(dstFile)
		if err != nil {
			t.Fatalf("file not copied: %s", err)
		}
		if fileInfo.Mode() != perm {
			t.Errorf("unexpected file permissions: got %v, want %v", fileInfo.Mode(), perm)
		}

		// read the contents of the copied file
		data, err := os.ReadFile(dstFile)
		if err != nil {
			t.Fatal(err)
		}

		// verify the contents of the copied file
		expected := []byte("test data")
		if string(data) != string(expected) {
			t.Errorf("unexpected file content: got %q, want %q", data, expected)
		}
	})

	t.Run("skip existing", func(t *testing.T) {
		// existing file in the destination dir
		path := filepath.Join(dstDir, "existing.txt")
		err := os.WriteFile(path, []byte("v1"), 0444)
		if err != nil {
			t.Fatal(err)
		}

		// same file in the source dir
		path = filepath.Join(srcDir, "existing.txt")
		err = os.WriteFile(path, []byte("v2"), 0444)
		if err != nil {
			t.Fatal(err)
		}

		// copy files
		pattern := filepath.Join(srcDir, "*.txt")
		err = CopyFiles(pattern, dstDir, 0444)
		if err != nil {
			t.Fatal(err)
		}

		// verify that the new file was copied correctly
		newFile := filepath.Join(dstDir, "source.txt")
		_, err = os.Stat(newFile)
		if err != nil {
			t.Fatalf("new file not copied: %s", err)
		}

		// verify that the existing file remained unchanged
		existFile := filepath.Join(dstDir, "existing.txt")
		data, err := os.ReadFile(existFile)
		if err != nil {
			t.Fatal(err)
		}
		expected := []byte("v1")
		if string(data) != string(expected) {
			t.Error("existing file got overwritten")
		}
	})
}

func TestReadJson(t *testing.T) {
	type Person struct{ Name string }

	t.Run("valid", func(t *testing.T) {
		got, err := ReadJson[Person](filepath.Join("testdata", "valid.json"))
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		want := Person{"alice"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ReadJson[Person](filepath.Join("testdata", "invalid.json"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		_, err := ReadJson[Person](filepath.Join("testdata", "missing.json"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestWriteFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "files")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	t.Run("text", func(t *testing.T) {
		path := filepath.Join(dir, "data.txt")
		err = WriteFile(path, "hello", 0444)
		if err != nil {
			t.Fatalf("expected nil err, got %v", err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read file: expected nil err, got %v", err)
		}
		want := []byte("hello")
		if !reflect.DeepEqual(got, want) {
			t.Errorf("read file: expected %v, got %v", want, got)
		}
	})

	t.Run("data-octet-stream", func(t *testing.T) {
		path := filepath.Join(dir, "data-1.bin")
		err = WriteFile(path, "data:application/octet-stream;base64,MTIz", 0444)
		if err != nil {
			t.Fatalf("expected nil err, got %v", err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read file: expected nil err, got %v", err)
		}
		want := []byte("123")
		if !reflect.DeepEqual(got, want) {
			t.Errorf("read file: expected %v, got %v", want, got)
		}
	})

	t.Run("data-base64", func(t *testing.T) {
		path := filepath.Join(dir, "data-2.bin")
		err = WriteFile(path, "data:;base64,MTIz", 0444)
		if err != nil {
			t.Fatalf("expected nil err, got %v", err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read file: expected nil err, got %v", err)
		}
		want := []byte("123")
		if !reflect.DeepEqual(got, want) {
			t.Errorf("read file: expected %v, got %v", want, got)
		}
	})

	t.Run("data-text-plain", func(t *testing.T) {
		path := filepath.Join(dir, "data-3.bin")
		err = WriteFile(path, "data:text/plain;,123", 0444)
		if err != nil {
			t.Fatalf("expected nil err, got %v", err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read file: expected nil err, got %v", err)
		}
		want := []byte("123")
		if !reflect.DeepEqual(got, want) {
			t.Errorf("read file: expected %v, got %v", want, got)
		}
	})

	t.Run("perm", func(t *testing.T) {
		const perm = 0444
		path := filepath.Join(dir, "perm.txt")
		err = WriteFile(path, "hello", perm)
		if err != nil {
			t.Fatalf("expected nil err, got %v", err)
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			t.Fatalf("file not created: %s", err)
		}
		if fileInfo.Mode().Perm() != perm {
			t.Errorf("unexpected file permissions: expected %o, got %o", perm, fileInfo.Mode().Perm())
		}
	})

	t.Run("missing data-url separator", func(t *testing.T) {
		path := filepath.Join(dir, "data.bin")
		err = WriteFile(path, "data:application/octet-stream:MTIz", 0444)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid binary value", func(t *testing.T) {
		path := filepath.Join(dir, "data.bin")
		err = WriteFile(path, "data:;base64,12345", 0444)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
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
			if (err != nil) != tt.wantErr {
				t.Errorf("JoinDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JoinDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMkdirTemp(t *testing.T) {
	t.Run("default permissions", func(t *testing.T) {
		const perm = 0755
		dir, err := MkdirTemp(perm)
		if err != nil {
			t.Fatalf("failed to create temp directory: %v", err)
		}
		defer os.Remove(dir)

		info, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("failed to stat temp directory: %v", err)
		}
		if info.Mode().Perm() != perm {
			t.Errorf("unexpected permissions: expected %o, got %o", perm, info.Mode().Perm())
		}
	})

	t.Run("non-default permissions", func(t *testing.T) {
		const perm = 0777
		dir, err := MkdirTemp(perm)
		if err != nil {
			t.Fatalf("failed to create temp directory: %v", err)
		}
		defer os.Remove(dir)

		info, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("failed to stat temp directory: %v", err)
		}
		if info.Mode().Perm() != perm {
			t.Errorf("unexpected permissions: expected %o, got %o", perm, info.Mode().Perm())
		}
	})
}
