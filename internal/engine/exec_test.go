package engine

import (
	"errors"
	"strings"
	"testing"

	"github.com/nalgeon/be"
	"github.com/nalgeon/codapi/internal/execy"
)

func TestProgram_Run(t *testing.T) {
	commands := map[string]execy.CmdOut{
		"mock stdout": {Stdout: "stdout", Stderr: "", Err: nil},
		"mock stderr": {Stdout: "", Stderr: "stderr", Err: nil},
		"mock outerr": {Stdout: "stdout", Stderr: "stderr", Err: nil},
		"mock err":    {Stdout: "", Stderr: "stderr", Err: errors.New("error")},
	}
	mem := execy.Mock(commands)

	for key, want := range commands {
		t.Run(key, func(t *testing.T) {
			p := NewProgram(3, 100)
			name, arg, _ := strings.Cut(key, " ")
			stdout, stderr, err := p.Run("mock_42", name, arg)
			be.True(t, mem.Has(key))
			be.Equal(t, stdout, want.Stdout)
			be.Equal(t, stderr, want.Stderr)
			be.Err(t, err, want.Err)
		})
	}
}

func TestProgram_LimitOutput(t *testing.T) {
	commands := map[string]execy.CmdOut{
		"mock stdout": {Stdout: "1234567890", Stderr: ""},
		"mock stderr": {Stdout: "", Stderr: "1234567890"},
		"mock outerr": {Stdout: "1234567890", Stderr: "0987654321"},
	}
	execy.Mock(commands)

	const nOutput = 5
	{
		p := NewProgram(3, nOutput)
		stdout, _, _ := p.Run("mock_42", "mock", "stdout")
		be.Equal(t, stdout, "12345")
	}
	{
		p := NewProgram(3, nOutput)
		_, stderr, _ := p.Run("mock_42", "mock", "stderr")
		be.Equal(t, stderr, "12345")
	}
	{
		p := NewProgram(3, nOutput)
		stdout, stderr, _ := p.Run("mock_42", "mock", "outerr")
		be.Equal(t, stdout, "12345")
		be.Equal(t, stderr, "09876")
	}
}
