package sandbox

import (
	"testing"

	"github.com/nalgeon/be"
)

func TestSemaphore(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		sem := NewSemaphore(3)
		be.Equal(t, sem.Size(), 3)
	})
	t.Run("acquire", func(t *testing.T) {
		sem := NewSemaphore(2)
		err := sem.Acquire()
		be.Err(t, err, nil)
		err = sem.Acquire()
		be.Err(t, err, nil)
		err = sem.Acquire()
		be.Err(t, err, ErrBusy)
	})
	t.Run("release", func(t *testing.T) {
		sem := NewSemaphore(2)
		_ = sem.Acquire()
		_ = sem.Acquire()
		_ = sem.Acquire()

		sem.Release()
		err := sem.Acquire()
		be.Err(t, err, nil)
	})
	t.Run("release free", func(t *testing.T) {
		sem := NewSemaphore(2)
		sem.Release()
		sem.Release()
		sem.Release()
	})
}
