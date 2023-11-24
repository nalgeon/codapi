package execy

import (
	"os/exec"
	"strings"

	"github.com/nalgeon/codapi/logx"
)

// Mock installs mock outputs for given commands.
func Mock(commands map[string]CmdOut) *logx.Memory {
	if commands != nil {
		mockCommands = commands
	}
	mem := logx.NewMemory("exec")
	runner = &mockRunner{mem}
	return mem
}

// mockRunner returns mock outputs
// without running OS programs.
type mockRunner struct {
	mem *logx.Memory
}

// Run returns a mock output from the registry
// that matches the given command name and argument.
func (r *mockRunner) Run(cmd *exec.Cmd) error {
	cmdStr := strings.Join(cmd.Args, " ")
	r.mem.WriteString(cmdStr)

	key := cmd.Args[0] + " " + cmd.Args[1]
	out, ok := mockCommands[key]
	if !ok {
		// command is not in the registry,
		// so let's return an empty "success" result
		out = CmdOut{}
	}
	_, _ = cmd.Stdout.Write([]byte(out.Stdout))
	_, _ = cmd.Stderr.Write([]byte(out.Stderr))
	return out.Err
}

var mockCommands map[string]CmdOut = map[string]CmdOut{}
