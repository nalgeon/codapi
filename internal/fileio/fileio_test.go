package fileio

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCopyFiles(t *testing.T) {
	// Create a temporary directory for testing
	srcDir, err := os.MkdirTemp("", "src")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	// Create a source file
	srcFile := filepath.Join(srcDir, "source.txt")
	err = os.WriteFile(srcFile, []byte("test data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Specify the destination directory
	dstDir, err := os.MkdirTemp("", "dst")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dstDir)

	// Call the CopyFiles function
	pattern := filepath.Join(srcDir, "*.txt")
	err = CopyFiles(pattern, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the file was copied correctly
	dstFile := filepath.Join(dstDir, "source.txt")
	_, err = os.Stat(dstFile)
	if err != nil {
		t.Fatalf("file not copied: %s", err)
	}

	// Read the contents of the copied file
	data, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the contents of the copied file
	expected := []byte("test data")
	if string(data) != string(expected) {
		t.Errorf("unexpected file content: got %q, want %q", data, expected)
	}
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
