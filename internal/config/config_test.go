package config

import (
	"reflect"
	"strings"
	"testing"
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BoxNames: expected %v, got %v", want, got)
	}
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CommandNames: expected %v, got %v", want, got)
	}
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
	if !strings.Contains(got, `"pool_size": 8`) {
		t.Error("ToJSON: expected pool_size = 8")
	}
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
	if box.Image != "" {
		t.Error("Image: should not set default value")
	}
	if box.Runtime != defs.Runtime {
		t.Errorf("Runtime: expected %s, got %s", defs.Runtime, box.Runtime)
	}
	if box.CPU != defs.CPU {
		t.Errorf("CPU: expected %d, got %d", defs.CPU, box.CPU)
	}
	if box.Memory != defs.Memory {
		t.Errorf("Memory: expected %d, got %d", defs.Memory, box.Memory)
	}
	if box.Storage != defs.Storage {
		t.Errorf("Storage: expected %s, got %s", defs.Storage, box.Storage)
	}
	if box.Network != defs.Network {
		t.Errorf("Network: expected %s, got %s", defs.Network, box.Network)
	}
	if box.Volume != defs.Volume {
		t.Errorf("Volume: expected %s, got %s", defs.Volume, box.Volume)
	}
	if !reflect.DeepEqual(box.Tmpfs, defs.Tmpfs) {
		t.Errorf("Tmpfs: expected %v, got %v", defs.Tmpfs, box.Tmpfs)
	}
	if !reflect.DeepEqual(box.CapAdd, defs.CapAdd) {
		t.Errorf("CapAdd: expected %v, got %v", defs.CapAdd, box.CapAdd)
	}
	if !reflect.DeepEqual(box.CapDrop, defs.CapDrop) {
		t.Errorf("CapDrop: expected %v, got %v", defs.CapDrop, box.CapDrop)
	}
	if !reflect.DeepEqual(box.Ulimit, defs.Ulimit) {
		t.Errorf("Ulimit: expected %v, got %v", defs.Ulimit, box.Ulimit)
	}
	if box.NProc != defs.NProc {
		t.Errorf("NProc: expected %d, got %d", defs.NProc, box.NProc)
	}
	if len(box.Files) != 0 {
		t.Error("Files: should not set default value")
	}
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
	if step.Box != "" {
		t.Error("Box: should not set default value")
	}
	if step.User != defs.User {
		t.Errorf("User: expected %s, got %s", defs.User, step.User)
	}
	if step.Action != defs.Action {
		t.Errorf("Action: expected %s, got %s", defs.Action, step.Action)
	}
	if len(step.Command) != 0 {
		t.Error("Command: should not set default value")
	}
	if step.Timeout != defs.Timeout {
		t.Errorf("Timeout: expected %d, got %d", defs.Timeout, step.Timeout)
	}
	if step.NOutput != defs.NOutput {
		t.Errorf("NOutput: expected %d, got %d", defs.NOutput, step.NOutput)
	}
}
