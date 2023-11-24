// Package sandbox provides a registry of sandboxes
// for code execution.
package sandbox

import (
	"errors"
	"strings"
	"time"

	"github.com/nalgeon/codapi/engine"
)

var ErrUnknownSandbox = errors.New("unknown sandbox")
var ErrUnknownCommand = errors.New("unknown command")
var ErrEmptyRequest = errors.New("empty request")

// Validate checks if the code execution request is valid.
func Validate(in engine.Request) error {
	box, ok := engines[in.Sandbox]
	if !ok {
		return ErrUnknownSandbox
	}
	_, ok = box[in.Command]
	if !ok {
		return ErrUnknownCommand
	}
	if len(in.Files) < 2 && strings.TrimSpace(in.Files.First()) == "" {
		return ErrEmptyRequest
	}
	return nil
}

// Exec executes the code using the appropriate sandbox.
// Allows no more than pool.Size() concurrent workers at any given time.
// The request must already be validated by Validate().
func Exec(in engine.Request) engine.Execution {
	err := semaphore.Acquire()
	defer semaphore.Release()
	if err == ErrBusy {
		return engine.Fail(in.ID, engine.ErrBusy)
	}
	start := time.Now()
	engine := engines[in.Sandbox][in.Command]
	out := engine.Exec(in)
	out.Duration = int(time.Since(start).Milliseconds())
	return out
}
