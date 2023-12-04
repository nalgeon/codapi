package engine

import (
	"errors"
	"reflect"
	"sort"
	"testing"
)

func TestExecutionError(t *testing.T) {
	inner := errors.New("inner error")
	err := NewExecutionError("failed", inner)
	if err.Error() != "failed: inner error" {
		t.Errorf("Error: expected %q, got %q", "failed: inner error", err.Error())
	}
	unwrapped := err.Unwrap()
	if unwrapped != inner {
		t.Errorf("Unwrap: expected %#v, got %#v", inner, unwrapped)
	}
}

func TestFiles_Count(t *testing.T) {
	var files Files = map[string]string{
		"first":  "alice",
		"second": "bob",
		"third":  "cindy",
	}
	if files.Count() != 3 {
		t.Errorf("Count: expected 3, got %d", files.Count())
	}
}

func TestFiles_Range(t *testing.T) {
	var files Files = map[string]string{
		"first":  "alice",
		"second": "bob",
		"third":  "cindy",
	}

	t.Run("range", func(t *testing.T) {
		names := []string{}
		contents := []string{}
		files.Range(func(name, content string) bool {
			names = append(names, name)
			contents = append(contents, content)
			return true
		})
		sort.Strings(names)
		if !reflect.DeepEqual(names, []string{"first", "second", "third"}) {
			t.Errorf("unexpected names: %v", names)
		}
		sort.Strings(contents)
		if !reflect.DeepEqual(contents, []string{"alice", "bob", "cindy"}) {
			t.Errorf("unexpected contents: %v", contents)
		}
	})

	t.Run("break", func(t *testing.T) {
		names := []string{}
		contents := []string{}
		files.Range(func(name, content string) bool {
			names = append(names, name)
			contents = append(contents, content)
			return false
		})
		if len(names) != 1 {
			t.Fatalf("expected names len = 1, got %d", len(names))
		}
		if len(contents) != 1 {
			t.Fatalf("expected contents len = 1, got %d", len(contents))
		}
		if files[names[0]] != contents[0] {
			t.Fatalf("name does not match content: %v -> %v", names[0], contents[0])
		}
	})
}

func TestFail(t *testing.T) {
	t.Run("ExecutionError", func(t *testing.T) {
		err := NewExecutionError("failed", errors.New("inner error"))
		out := Fail("42", err)
		if out.ID != "42" {
			t.Errorf("ID: expected 42, got %v", out.ID)
		}
		if out.OK {
			t.Error("OK: expected false")
		}
		if out.Stderr != "internal error" {
			t.Errorf("Stderr: expected %q, got %q", "internal error", out.Stderr)
		}
		if out.Stdout != "" {
			t.Errorf("Stdout: expected empty, got %q", out.Stdout)
		}
		if out.Err != err {
			t.Errorf("Err: expected %#v, got %#v", err, out.Err)
		}
	})
	t.Run("ErrBusy", func(t *testing.T) {
		err := ErrBusy
		out := Fail("42", err)
		if out.ID != "42" {
			t.Errorf("ID: expected 42, got %v", out.ID)
		}
		if out.OK {
			t.Error("OK: expected false")
		}
		if out.Stderr != err.Error() {
			t.Errorf("Stderr: expected %q, got %q", err.Error(), out.Stderr)
		}
		if out.Stdout != "" {
			t.Errorf("Stdout: expected empty, got %q", out.Stdout)
		}
		if out.Err != err {
			t.Errorf("Err: expected %#v, got %#v", err, out.Err)
		}
	})
	t.Run("Error", func(t *testing.T) {
		err := errors.New("user error")
		out := Fail("42", err)
		if out.ID != "42" {
			t.Errorf("ID: expected 42, got %v", out.ID)
		}
		if out.OK {
			t.Error("OK: expected false")
		}
		if out.Stderr != err.Error() {
			t.Errorf("Stderr: expected %q, got %q", err.Error(), out.Stderr)
		}
		if out.Stdout != "" {
			t.Errorf("Stdout: expected empty, got %q", out.Stdout)
		}
		if out.Err != nil {
			t.Errorf("Err: expected nil, got %#v", out.Err)
		}
	})

}
