package sandbox

import (
	"testing"

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
	if err != nil {
		t.Fatalf("ApplyConfig: expected nil err, got %v", err)
	}
	if semaphore.Size() != cfg.PoolSize {
		t.Errorf("semaphore.Size: expected %d, got %d", cfg.PoolSize, semaphore.Size())
	}
	if len(engines) != 2 {
		t.Errorf("len(engines): expected 2, got %d", len(engines))
	}
	if len(engines["http"]) != 1 {
		t.Errorf("len(engine = http): expected 1, got %d", len(engines["http"]))
	}
	if _, ok := engines["http"]["run"].(*engine.HTTP); !ok {
		t.Error("engine = http: expected HTTP engine")
	}
	if len(engines["python"]) != 2 {
		t.Errorf("len(engine = python): expected 2, got %d", len(engines["python"]))
	}
	if _, ok := engines["python"]["run"].(*engine.Docker); !ok {
		t.Error("engine = python: expected Docker engine")
	}
}
