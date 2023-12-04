package engine

import (
	"strings"
	"testing"

	"github.com/nalgeon/codapi/internal/config"
	"github.com/nalgeon/codapi/internal/execy"
	"github.com/nalgeon/codapi/internal/logx"
)

var dockerCfg = &config.Config{
	Boxes: map[string]*config.Box{
		"postgresql": {
			Image:   "codapi/postgresql",
			Runtime: "runc",
			Host: config.Host{
				CPU: 1, Memory: 64, Network: "none",
				Volume: "%s:/sandbox:ro",
				NProc:  64,
			},
		},
		"python": {
			Image:   "codapi/python",
			Runtime: "runc",
			Host: config.Host{
				CPU: 1, Memory: 64, Network: "none",
				Volume: "%s:/sandbox:ro",
				NProc:  64,
			},
		},
	},
	Commands: map[string]config.SandboxCommands{
		"postgresql": map[string]*config.Command{
			"run": {
				Engine: "docker",
				Before: &config.Step{
					Box: "postgres", User: "sandbox", Action: "exec",
					Command: []string{"psql", "-f", "create.sql"},
					NOutput: 4096,
				},
				Steps: []*config.Step{
					{
						Box: "postgres", User: "sandbox", Action: "exec", Stdin: true,
						Command: []string{"psql", "--user=:name"},
						NOutput: 4096,
					},
				},
				After: &config.Step{
					Box: "postgres", User: "sandbox", Action: "exec",
					Command: []string{"psql", "-f", "drop.sql"},
					NOutput: 4096,
				},
			},
		},
		"python": map[string]*config.Command{
			"run": {
				Engine: "docker",
				Entry:  "main.py",
				Steps: []*config.Step{
					{
						Box: "python", User: "sandbox", Action: "run",
						Command: []string{"python", "main.py"},
						NOutput: 4096,
					},
				},
			},
		},
	},
}

func TestDockerRun(t *testing.T) {
	logx.Mock()
	commands := map[string]execy.CmdOut{
		"docker run": {Stdout: "hello world", Stderr: "", Err: nil},
	}
	mem := execy.Mock(commands)
	engine := NewDocker(dockerCfg, "python", "run")

	t.Run("success", func(t *testing.T) {
		req := Request{
			ID:      "http_42",
			Sandbox: "python",
			Command: "run",
			Files: map[string]string{
				"": "print('hello world')",
			},
		}
		out := engine.Exec(req)
		if out.ID != req.ID {
			t.Errorf("ID: expected %s, got %s", req.ID, out.ID)
		}
		if !out.OK {
			t.Error("OK: expected true")
		}
		want := "hello world"
		if out.Stdout != want {
			t.Errorf("Stdout: expected %q, got %q", want, out.Stdout)
		}
		if out.Stderr != "" {
			t.Errorf("Stderr: expected %q, got %q", "", out.Stdout)
		}
		if out.Err != nil {
			t.Errorf("Err: expected nil, got %v", out.Err)
		}
		mem.MustHave(t, "python main.py")
	})
}

func TestDockerExec(t *testing.T) {
	logx.Mock()
	commands := map[string]execy.CmdOut{
		"docker exec": {Stdout: "hello world", Stderr: "", Err: nil},
	}
	mem := execy.Mock(commands)
	engine := NewDocker(dockerCfg, "postgresql", "run")

	t.Run("success", func(t *testing.T) {
		req := Request{
			ID:      "http_42",
			Sandbox: "postgresql",
			Command: "run",
			Files: map[string]string{
				"": "select 'hello world'",
			},
		}
		out := engine.Exec(req)
		if out.ID != req.ID {
			t.Errorf("ID: expected %s, got %s", req.ID, out.ID)
		}
		if !out.OK {
			t.Error("OK: expected true")
		}
		want := "hello world"
		if out.Stdout != want {
			t.Errorf("Stdout: expected %q, got %q", want, out.Stdout)
		}
		if out.Stderr != "" {
			t.Errorf("Stderr: expected %q, got %q", "", out.Stdout)
		}
		if out.Err != nil {
			t.Errorf("Err: expected nil, got %v", out.Err)
		}
		mem.MustHave(t, "psql -f create.sql")
		mem.MustHave(t, "psql --user=http_42")
		mem.MustHave(t, "psql -f drop.sql")
	})
}

func Test_expandVars(t *testing.T) {
	const name = "codapi_01"
	commands := map[string]string{
		"python main.py":     "python main.py",
		"sh create.sh :name": "sh create.sh " + name,
	}
	for cmd, want := range commands {
		src := strings.Fields(cmd)
		exp := expandVars(src, name)
		got := strings.Join(exp, " ")
		if got != want {
			t.Errorf("%q: expected %q, got %q", cmd, got, want)
		}
	}
}
