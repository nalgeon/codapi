package execy

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/nalgeon/be"
)

func TestRunner(t *testing.T) {
	const want = "hello world"
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "echo", "-n", want)
	outb := new(strings.Builder)
	errb := new(strings.Builder)
	cmd.Stdout = outb
	cmd.Stderr = errb

	err := Run(cmd)
	be.Err(t, err, nil)
	be.Equal(t, outb.String(), want)
	be.Equal(t, errb.String(), "")
}
