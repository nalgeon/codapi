package engine

import (
	"errors"
	"strings"
	"testing"

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
			if !mem.Has(key) {
				t.Errorf("Run: command %q not run", key)
			}
			if stdout != want.Stdout {
				t.Errorf("stdout: want %#v, got %#v", want.Stdout, stdout)
			}
			if stderr != want.Stderr {
				t.Errorf("stderr: want %#v, got %#v", want.Stderr, stderr)
			}
			if err != want.Err {
				t.Errorf("err: want %#v, got %#v", want.Err, err)
			}
		})
	}
}

func TestProgram_LimitOutput(t *testing.T) {
	commands := map[string]execy.CmdOut{
		"mock stdout": {Stdout: "1234567890", Stderr: ""},
		"mock stderr": {Stdout: "", Stderr: "1234567890"},
		"mock outerr": {Stdout: "1234567890", Stderr: "1234567890"},
	}
	execy.Mock(commands)

	const nOutput = 5
	{
		p := NewProgram(3, nOutput)
		stdout, _, _ := p.Run("mock_42", "mock", "stdout")
		if stdout != "12345" {
			t.Errorf("stdout: want %#v, got %#v", "12345", stdout)
		}
	}
	{
		p := NewProgram(3, nOutput)
		_, stderr, _ := p.Run("mock_42", "mock", "stderr")
		if stderr != "12345" {
			t.Errorf("stderr: want %#v, got %#v", "12345", stderr)
		}
	}
	{
		p := NewProgram(3, nOutput)
		stdout, stderr, _ := p.Run("mock_42", "mock", "outerr")
		if stdout != "12345" {
			t.Errorf("stdout: want %#v, got %#v", "12345", stdout)
		}
		if stderr != "12345" {
			t.Errorf("stderr: want %#v, got %#v", "12345", stderr)
		}
	}
}
