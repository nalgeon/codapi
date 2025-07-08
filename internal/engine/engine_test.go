package engine

import (
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/nalgeon/be"
)

func TestGenerateID(t *testing.T) {
	t.Run("with version", func(t *testing.T) {
		req := Request{
			Sandbox: "python",
			Version: "dev",
			Command: "run",
		}
		req.GenerateID()
		be.True(t, strings.HasPrefix(req.ID, "python.dev_run_"))
	})
	t.Run("without version", func(t *testing.T) {
		req := Request{
			Sandbox: "python",
			Command: "run",
		}
		req.GenerateID()
		be.True(t, strings.HasPrefix(req.ID, "python_run_"))
	})
}

func TestExecutionError(t *testing.T) {
	inner := errors.New("inner error")
	err := NewExecutionError("failed", inner)
	be.Err(t, err, inner)
}

func TestFiles_Count(t *testing.T) {
	var files Files = map[string]string{
		"first":  "alice",
		"second": "bob",
		"third":  "cindy",
	}
	be.Equal(t, files.Count(), 3)
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
		be.Equal(t, names, []string{"first", "second", "third"})
		sort.Strings(contents)
		be.Equal(t, contents, []string{"alice", "bob", "cindy"})
	})

	t.Run("break", func(t *testing.T) {
		names := []string{}
		contents := []string{}
		files.Range(func(name, content string) bool {
			names = append(names, name)
			contents = append(contents, content)
			return false
		})
		be.Equal(t, len(names), 1)
		be.Equal(t, len(contents), 1)
		be.Equal(t, files[names[0]], contents[0])
	})
}

func TestFail(t *testing.T) {
	t.Run("ExecutionError", func(t *testing.T) {
		err := NewExecutionError("failed", errors.New("inner error"))
		out := Fail("42", err)
		be.Equal(t, out.ID, "42")
		be.Equal(t, out.OK, false)
		be.Equal(t, out.Stderr, "internal error")
		be.Equal(t, out.Stdout, "")
		be.Err(t, out.Err, err)
	})
	t.Run("ErrBusy", func(t *testing.T) {
		err := ErrBusy
		out := Fail("42", err)
		be.Equal(t, out.ID, "42")
		be.Equal(t, out.OK, false)
		be.Equal(t, out.Stderr, err.Error())
		be.Equal(t, out.Stdout, "")
		be.Err(t, out.Err, err)
	})
	t.Run("Error", func(t *testing.T) {
		err := errors.New("user error")
		out := Fail("42", err)
		be.Equal(t, out.ID, "42")
		be.Equal(t, out.OK, false)
		be.Equal(t, out.Stderr, err.Error())
		be.Equal(t, out.Stdout, "")
		be.Err(t, out.Err, nil)
	})
}
