// Package execy runs external commands.
package execy

import (
	"os/exec"
)

var runner = Runner(&osRunner{})

// Runner executes external commands.
type Runner interface {
	Run(cmd *exec.Cmd) error
}

// osRunner runs OS programs.
type osRunner struct{}

func (r *osRunner) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}

func Run(cmd *exec.Cmd) error {
	return runner.Run(cmd)
}

// CmdOut represents the result of the command run.
type CmdOut struct {
	Stdout string
	Stderr string
	Err    error
}
