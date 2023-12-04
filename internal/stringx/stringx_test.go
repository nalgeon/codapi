package stringx

import "testing"

func TestShorten(t *testing.T) {
	t.Run("shorten", func(t *testing.T) {
		const src = "Hello, World!"
		const want = "Hello [truncated]"
		got := Shorten(src, 5)
		if got != want {
			t.Errorf("expected %q, got %q", got, want)
		}
	})
	t.Run("ignore", func(t *testing.T) {
		const src = "Hello, World!"
		const want = src
		got := Shorten(src, 20)
		if got != want {
			t.Errorf("expected %q, got %q", got, want)
		}
	})
}

func TestCompact(t *testing.T) {
	t.Run("compact", func(t *testing.T) {
		const src = "go\nis   awesome"
		const want = "go is awesome"
		got := Compact(src)
		if got != want {
			t.Errorf("expected %q, got %q", got, want)
		}
	})
	t.Run("ignore", func(t *testing.T) {
		const src = "go is awesome"
		const want = src
		got := Compact(src)
		if got != want {
			t.Errorf("expected %q, got %q", got, want)
		}
	})
}

func TestRandString(t *testing.T) {
	lengths := []int{2, 4, 6, 8, 10}
	for _, n := range lengths {
		s := RandString(n)
		if len(s) != n {
			t.Errorf("%d: expected len(s) = %d, got %d", n, n, len(s))
		}
	}
}
