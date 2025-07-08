package logx

import (
	"testing"

	"github.com/nalgeon/be"
)

func TestMemory_Name(t *testing.T) {
	mem := NewMemory("log")
	be.Equal(t, mem.Name, "log")
}

func TestMemory_Write(t *testing.T) {
	mem := NewMemory("log")
	be.Equal(t, len(mem.Lines), 0)

	n, err := mem.Write([]byte("hello world"))
	be.Err(t, err, nil)
	be.Equal(t, n, 11)

	be.Equal(t, len(mem.Lines), 1)
	be.Equal(t, mem.Lines[0], "hello world")
}

func TestMemory_Has(t *testing.T) {
	mem := NewMemory("log")
	be.True(t, !mem.Has("hello world"))
	_, _ = mem.Write([]byte("hello world"))
	be.True(t, mem.Has("hello world"))
	_, _ = mem.Write([]byte("one two three four"))
	be.True(t, mem.Has("one two"))
	be.True(t, mem.Has("two three"))
	be.True(t, !mem.Has("one three"))
	be.True(t, mem.Has("one", "three"))
	be.True(t, mem.Has("one", "three", "four"))
	be.True(t, mem.Has("four", "three"))
}
