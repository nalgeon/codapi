package sandbox

import (
	"testing"

	"github.com/nalgeon/be"
	"github.com/nalgeon/codapi/internal/engine"
	"github.com/nalgeon/codapi/internal/execy"
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
		be.Err(t, err, nil)
	})
	t.Run("unknown sandbox", func(t *testing.T) {
		req := engine.Request{
			ID:      "http_42",
			Sandbox: "rust",
			Command: "run",
			Files:   nil,
		}
		err := Validate(req)
		be.Err(t, err, ErrUnknownSandbox)
	})
	t.Run("unknown command", func(t *testing.T) {
		req := engine.Request{
			ID:      "http_42",
			Sandbox: "python",
			Command: "deploy",
			Files:   nil,
		}
		err := Validate(req)
		be.Err(t, err, ErrUnknownCommand)
	})
	t.Run("empty request", func(t *testing.T) {
		req := engine.Request{
			ID:      "http_42",
			Sandbox: "python",
			Command: "run",
			Files:   nil,
		}
		err := Validate(req)
		be.Err(t, err, ErrEmptyRequest)
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
		be.Equal(t, out.ID, req.ID)
		be.True(t, out.OK)
		be.Equal(t, out.Stdout, "hello")
		be.Equal(t, out.Stderr, "")
		be.Equal(t, out.Err, nil)
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
		be.Err(t, out.Err, engine.ErrBusy)
	})
}
