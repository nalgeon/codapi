package sandbox

import "testing"

func TestSemaphore(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		sem := NewSemaphore(3)
		if sem.Size() != 3 {
			t.Errorf("Size: expected 3, got %d", sem.Size())
		}
	})
	t.Run("acquire", func(t *testing.T) {
		sem := NewSemaphore(2)
		err := sem.Acquire()
		if err != nil {
			t.Fatalf("acquire #1: expected nil err")
		}
		err = sem.Acquire()
		if err != nil {
			t.Fatalf("acquire #2: expected nil err")
		}
		err = sem.Acquire()
		if err != ErrBusy {
			t.Fatalf("acquire #3: expected ErrBusy")
		}
	})
	t.Run("release", func(t *testing.T) {
		sem := NewSemaphore(2)
		_ = sem.Acquire()
		_ = sem.Acquire()
		_ = sem.Acquire()

		sem.Release()
		err := sem.Acquire()
		if err != nil {
			t.Fatalf("acquire after release: expected nil err")
		}
	})
	t.Run("release free", func(t *testing.T) {
		sem := NewSemaphore(2)
		sem.Release()
		sem.Release()
		sem.Release()
	})
}
