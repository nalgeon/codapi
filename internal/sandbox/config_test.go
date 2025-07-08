package sandbox

import (
	"testing"

	"github.com/nalgeon/be"
	"github.com/nalgeon/codapi/internal/config"
	"github.com/nalgeon/codapi/internal/engine"
)

var cfg = &config.Config{
	PoolSize: 8,
	HTTP: &config.HTTP{
		Hosts: map[string]string{"localhost": "localhost"},
	},
	Boxes: map[string]*config.Box{
		"http":   {},
		"python": {},
	},
	Commands: map[string]config.SandboxCommands{
		"http": map[string]*config.Command{
			"run": {Engine: "http"},
		},
		"python": map[string]*config.Command{
			"run": {
				Engine: "docker",
				Entry:  "main.py",
				Steps: []*config.Step{
					{Box: "python", Action: "run", NOutput: 4096},
				},
			},
			"test": {Engine: "docker"},
		},
	},
}

func TestApplyConfig(t *testing.T) {
	err := ApplyConfig(cfg)
	be.Err(t, err, nil)
	be.Equal(t, semaphore.Size(), cfg.PoolSize)
	be.Equal(t, len(engines), 2)
	be.Equal(t, len(engines["http"]), 1)
	_, ok := engines["http"]["run"].(*engine.HTTP)
	be.True(t, ok)
	be.Equal(t, len(engines["python"]), 2)
	_, ok = engines["python"]["run"].(*engine.Docker)
	be.True(t, ok)
}
