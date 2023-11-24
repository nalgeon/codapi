package sandbox

import "errors"

var ErrBusy = errors.New("busy")

// A Semaphore manages a limited number of tokens
// that can be acquired or released.
type Semaphore struct {
	tokens chan struct{}
}

// NewSemaphore creates a new semaphore of the specified size.
func NewSemaphore(size int) *Semaphore {
	tokens := make(chan struct{}, size)
	for i := 0; i < size; i++ {
		tokens <- struct{}{}
	}
	return &Semaphore{tokens}
}

// Acquire acquires a token. Returns ErrBusy if no tokens are available.
func (q *Semaphore) Acquire() error {
	select {
	case <-q.tokens:
		return nil
	default:
		return ErrBusy
	}
}

// Release releases a token.
func (q *Semaphore) Release() {
	select {
	case q.tokens <- struct{}{}:
		return
	default:
		return
	}
}

// Size returns the size of the semaphore.
func (q *Semaphore) Size() int {
	return len(q.tokens)
}
