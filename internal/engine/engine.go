// Package engine provides code execution engines.
package engine

import (
	"errors"
	"fmt"

	"github.com/nalgeon/codapi/internal/stringx"
)

// A Request initiates code execution.
type Request struct {
	ID      string `json:"id"`
	Sandbox string `json:"sandbox"`
	Version string `json:"version,omitempty"`
	Command string `json:"command"`
	Files   Files  `json:"files"`
}

// GenerateID() sets a unique ID for the request.
func (r *Request) GenerateID() {
	if r.Version != "" {
		r.ID = fmt.Sprintf("%s.%s_%s_%s", r.Sandbox, r.Version, r.Command, stringx.RandString(8))
	} else {
		r.ID = fmt.Sprintf("%s_%s_%s", r.Sandbox, r.Command, stringx.RandString(8))
	}
}

// An Execution is an output from the code execution engine.
type Execution struct {
	ID       string `json:"id"`
	OK       bool   `json:"ok"`
	Duration int    `json:"duration"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Err      error  `json:"-"`
}

// An ErrTimeout is returned if code execution did not complete
// in the allowed timeframe.
var ErrTimeout = errors.New("code execution timeout")

// An ErrBusy is returned when there are no engines available.
var ErrBusy = errors.New("busy: try again later")

// An ExecutionError is returned if code execution failed
// due to the application problems, not due to the problems with the code.
type ExecutionError struct {
	msg   string
	inner error
}

func NewExecutionError(msg string, err error) ExecutionError {
	return ExecutionError{msg: msg, inner: err}
}

func (err ExecutionError) Error() string {
	return err.msg + ": " + err.inner.Error()
}

func (err ExecutionError) Unwrap() error {
	return err.inner
}

// An ArgumentError is returned if code execution failed
// due to the invalid value of the request argument.
type ArgumentError struct {
	name   string
	reason error
}

func NewArgumentError(name string, reason error) ArgumentError {
	return ArgumentError{name: name, reason: reason}
}

func (err ArgumentError) Error() string {
	return err.name + ": " + err.reason.Error()
}

func (err ArgumentError) Unwrap() error {
	return err.reason
}

// Files are a collection of files to be executed by the engine.
type Files map[string]string

// First returns the contents of the first file.
func (f Files) First() string {
	for _, content := range f {
		return content
	}
	return ""
}

// Range iterates over the files, calling fn for each one.
func (f Files) Range(fn func(name, content string) bool) {
	for name, content := range f {
		ok := fn(name, content)
		if !ok {
			break
		}
	}
}

// Count returns the number of files.
func (f Files) Count() int {
	return len(f)
}

// An Engine executes a specific sandbox command on the code.
// Engines must be concurrent-safe, since they can be accessed by multiple goroutines.
type Engine interface {
	// Exec executes the command and returns the output.
	Exec(req Request) Execution
}

// Fail creates an output from an error.
func Fail(id string, err error) Execution {
	if _, ok := err.(ExecutionError); ok {
		return Execution{
			ID:     id,
			OK:     false,
			Stderr: "internal error",
			Err:    err,
		}
	}
	if errors.Is(err, ErrBusy) {
		return Execution{
			ID:     id,
			OK:     false,
			Stderr: err.Error(),
			Err:    err,
		}
	}
	return Execution{
		ID:     id,
		OK:     false,
		Stderr: err.Error(),
	}
}
