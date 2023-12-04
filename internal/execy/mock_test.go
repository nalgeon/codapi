package execy

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)

func TestMock(t *testing.T) {
	const want = "hello world"
	out := CmdOut{Stdout: want, Stderr: "", Err: nil}
	mem := Mock(map[string]CmdOut{"echo -n": out})

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "echo", "-n", want)
	outb := new(strings.Builder)
	errb := new(strings.Builder)
	cmd.Stdout = outb
	cmd.Stderr = errb

	err := Run(cmd)
	if err != nil {
		t.Fatalf("Err: expected nil, got %v", err)
	}
	if outb.String() != want {
		t.Errorf("Stdout: expected %q, got %q", want, outb.String())
	}
	if errb.String() != "" {
		t.Errorf("Stderr: expected %q, got %q", "", errb.String())
	}
	mem.MustHave(t, "echo -n hello world")
}
