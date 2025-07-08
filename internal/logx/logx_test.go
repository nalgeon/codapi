package logx

import (
	"testing"

	"github.com/nalgeon/be"
)

func TestSetOutput(t *testing.T) {
	mem := NewMemory("log")
	SetOutput(mem)
	Log("hello")
	be.True(t, mem.Has("hello"))
}

func TestLog(t *testing.T) {
	mem := NewMemory("log")
	SetOutput(mem)
	{
		Log("value: %d", 42)
		be.Equal(t, len(mem.Lines), 1)
		be.True(t, mem.Has("value: 42"))
	}
	{
		Log("value: %d", 84)
		be.Equal(t, len(mem.Lines), 2)
		be.True(t, mem.Has("value: 42"))
		be.True(t, mem.Has("value: 84"))
	}
}

func TestDebug(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		mem := NewMemory("log")
		SetOutput(mem)
		Verbose = true
		{
			Debug("value: %d", 42)
			be.Equal(t, len(mem.Lines), 1)
			be.True(t, mem.Has("value: 42"))
		}
		{
			Debug("value: %d", 84)
			be.Equal(t, len(mem.Lines), 2)
			be.True(t, mem.Has("value: 42"))
			be.True(t, mem.Has("value: 84"))
		}
	})
	t.Run("disabled", func(t *testing.T) {
		mem := NewMemory("log")
		SetOutput(mem)
		Verbose = false
		Debug("value: %d", 42)
		be.Equal(t, len(mem.Lines), 0)
	})
}
