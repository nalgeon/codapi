package engine

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/nalgeon/codapi/internal/execy"
	"github.com/nalgeon/codapi/internal/logx"
)

// A Program is an executable program.
type Program struct {
	timeout time.Duration
	nOutput int64
}

// NewProgram creates a new program.
func NewProgram(timeoutSec int, nOutput int64) *Program {
	return &Program{
		timeout: time.Duration(timeoutSec) * time.Second,
		nOutput: nOutput,
	}
}

// Run starts the program and waits for it to complete (or timeout).
func (p *Program) Run(id, name string, arg ...string) (stdout string, stderr string, err error) {
	return p.RunStdin(nil, id, name, arg...)
}

// RunStdin starts the program with data from stdin
// and waits for it to complete (or timeout).
func (p *Program) RunStdin(stdin io.Reader, id, name string, arg ...string) (stdout string, stderr string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	var cmdout, cmderr strings.Builder
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Cancel = func() error {
		err := cmd.Process.Kill()
		logx.Debug("%s: execution timeout, killed process=%d, err=%v", id, cmd.Process.Pid, err)
		return err
	}

	cmd.Stdin = stdin
	cmd.Stdout = LimitWriter(&cmdout, p.nOutput)
	cmd.Stderr = LimitWriter(&cmderr, p.nOutput)
	err = execy.Run(cmd)
	stdout = strings.TrimSpace(cmdout.String())
	stderr = strings.TrimSpace(cmderr.String())
	return
}
