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
		// the worker returns its token when it is done
		be.Equal(t, semaphore.Size(), cfg.PoolSize)
	})
	t.Run("busy", func(t *testing.T) {
		defer func() { _ = ApplyConfig(cfg) }()
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
		// a rejected request must not release a token it never acquired
		be.Equal(t, semaphore.Size(), 0)
		// otherwise the next request picks up the donated token
		// and runs while the pool is still fully occupied
		out = Exec(req)
		be.Err(t, out.Err, engine.ErrBusy)
	})
}
