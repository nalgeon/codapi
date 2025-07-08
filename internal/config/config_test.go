package config

import (
	"strings"
	"testing"

	"github.com/nalgeon/be"
)

func TestConfig_BoxNames(t *testing.T) {
	cfg := &Config{
		Boxes: map[string]*Box{
			"go":     {},
			"python": {},
		},
	}

	want := []string{"go", "python"}
	got := cfg.BoxNames()
	be.Equal(t, got, want)
}

func TestConfig_CommandNames(t *testing.T) {
	cfg := &Config{
		Commands: map[string]SandboxCommands{
			"go": map[string]*Command{
				"run": {},
			},
			"python": map[string]*Command{
				"run":  {},
				"test": {},
			},
		},
	}

	want := []string{"go", "python"}
	got := cfg.CommandNames()
	be.Equal(t, got, want)
}

func TestConfig_ToJSON(t *testing.T) {
	cfg := &Config{
		PoolSize: 8,
		Boxes: map[string]*Box{
			"go":     {},
			"python": {},
		},
		Commands: map[string]SandboxCommands{
			"go": map[string]*Command{
				"run": {},
			},
			"python": map[string]*Command{
				"run":  {},
				"test": {},
			},
		},
	}

	got := cfg.ToJSON()
	be.True(t, strings.Contains(got, `"pool_size": 8`))
}

func Test_setBoxDefaults(t *testing.T) {
	box := &Box{}
	defs := &Box{
		Image:   "codapi/python",
		Runtime: "runc",
		Host: Host{
			CPU: 1, Memory: 64, Storage: "16m",
			Network: "none", Writable: true,
			Volume:  "%s:/sandbox:ro",
			Tmpfs:   []string{"/tmp:rw,size=16m"},
			CapAdd:  []string{"all"},
			CapDrop: []string{"none"},
			Ulimit:  []string{"nofile=96"},
			NProc:   96,
		},
		Files: []string{"config.py"},
	}
	setBoxDefaults(box, defs)
	be.Equal(t, box.Image, "")
	be.Equal(t, box.Runtime, defs.Runtime)
	be.Equal(t, box.CPU, defs.CPU)
	be.Equal(t, box.Memory, defs.Memory)
	be.Equal(t, box.Storage, defs.Storage)
	be.Equal(t, box.Network, defs.Network)
	be.Equal(t, box.Volume, defs.Volume)
	be.Equal(t, box.Tmpfs, defs.Tmpfs)
	be.Equal(t, box.CapAdd, defs.CapAdd)
	be.Equal(t, box.CapDrop, defs.CapDrop)
	be.Equal(t, box.Ulimit, defs.Ulimit)
	be.Equal(t, box.NProc, defs.NProc)
	be.Equal(t, len(box.Files), 0)
}

func Test_setStepDefaults(t *testing.T) {
	step := &Step{}
	defs := &Step{
		Box:     "python",
		User:    "sandbox",
		Action:  "run",
		Command: []string{"python", "main.py"},
		Timeout: 3,
		NOutput: 4096,
	}

	setStepDefaults(step, defs)
	be.Equal(t, step.Box, "")
	be.Equal(t, step.User, defs.User)
	be.Equal(t, step.Action, defs.Action)
	be.Equal(t, len(step.Command), 0)
	be.Equal(t, step.Timeout, defs.Timeout)
	be.Equal(t, step.NOutput, defs.NOutput)
}
