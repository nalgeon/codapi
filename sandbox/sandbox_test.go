package sandbox

import (
	"errors"
	"testing"

	"github.com/nalgeon/codapi/engine"
	"github.com/nalgeon/codapi/execy"
)

func TestValidate(t *testing.T) {
	_ = ApplyConfig(cfg)
	t.Run("valid", func(t *testing.T) {

		req := engine.Request{
			ID:      "http_42",
			Sandbox: "python",
			Command: "run",
			Files: map[string]string{
				"": "print('hello')",
			},
		}
		err := Validate(req)
		if err != nil {
			t.Errorf("Validate: expected nil err, got %v", err)
		}
	})
	t.Run("unknown sandbox", func(t *testing.T) {
		req := engine.Request{
			ID:      "http_42",
			Sandbox: "rust",
			Command: "run",
			Files:   nil,
		}
		err := Validate(req)
		if !errors.Is(err, ErrUnknownSandbox) {
			t.Errorf("Validate: expected ErrUnknownSandbox, got %T(%s)", err, err)
		}
	})
	t.Run("unknown command", func(t *testing.T) {
		req := engine.Request{
			ID:      "http_42",
			Sandbox: "python",
			Command: "deploy",
			Files:   nil,
		}
		err := Validate(req)
		if !errors.Is(err, ErrUnknownCommand) {
			t.Errorf("Validate: expected ErrUnknownCommand, got %T(%s)", err, err)
		}
	})
	t.Run("empty request", func(t *testing.T) {
		req := engine.Request{
			ID:      "http_42",
			Sandbox: "python",
			Command: "run",
			Files:   nil,
		}
		err := Validate(req)
		if !errors.Is(err, ErrEmptyRequest) {
			t.Errorf("Validate: expected ErrEmptyRequest, got %T(%s)", err, err)
		}
	})
}

func TestExec(t *testing.T) {
	_ = ApplyConfig(cfg)
	t.Run("exec", func(t *testing.T) {
		execy.Mock(map[string]execy.CmdOut{
			"docker run": {Stdout: "hello"},
		})
		req := engine.Request{
			ID:      "http_42",
			Sandbox: "python",
			Command: "run",
			Files: map[string]string{
				"": "print('hello')",
			},
		}
		out := Exec(req)
		if out.ID != req.ID {
			t.Errorf("ID: expected %s, got %s", req.ID, out.ID)
		}
		if !out.OK {
			t.Error("OK: expected true")
		}
		if out.Stdout != "hello" {
			t.Errorf("Stdout: expected hello, got %s", out.Stdout)
		}
		if out.Stderr != "" {
			t.Errorf("Stderr: expected empty string, got %s", out.Stderr)
		}
		if out.Err != nil {
			t.Errorf("Err: expected nil, got %v", out.Err)
		}
	})
	t.Run("busy", func(t *testing.T) {
		for i := 0; i < cfg.PoolSize; i++ {
			_ = semaphore.Acquire()
		}
		req := engine.Request{
			ID:      "http_42",
			Sandbox: "python",
			Command: "run",
			Files: map[string]string{
				"": "print('hello')",
			},
		}
		out := Exec(req)
		if out.Err != engine.ErrBusy {
			t.Errorf("Err: expected ErrBusy, got %v", out.Err)
		}
	})

}
