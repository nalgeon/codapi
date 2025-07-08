package stringx

import (
	"testing"

	"github.com/nalgeon/be"
)

func TestShorten(t *testing.T) {
	t.Run("shorten", func(t *testing.T) {
		const src = "Hello, World!"
		const want = "Hello [truncated]"
		got := Shorten(src, 5)
		be.Equal(t, got, want)
	})
	t.Run("ignore", func(t *testing.T) {
		const src = "Hello, World!"
		const want = src
		got := Shorten(src, 20)
		be.Equal(t, got, want)
	})
}

func TestCompact(t *testing.T) {
	t.Run("compact", func(t *testing.T) {
		const src = "go\nis   awesome"
		const want = "go is awesome"
		got := Compact(src)
		be.Equal(t, got, want)
	})
	t.Run("ignore", func(t *testing.T) {
		const src = "go is awesome"
		const want = src
		got := Compact(src)
		be.Equal(t, got, want)
	})
}

func TestRandString(t *testing.T) {
	lengths := []int{2, 4, 6, 8, 10}
	for _, n := range lengths {
		s := RandString(n)
		be.Equal(t, len(s), n)
	}
}
