package config

import (
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	cfgPath := filepath.Join("testdata", "config.json")
	boxPath := filepath.Join("testdata", "boxes.json")
	cmdPath := filepath.Join("testdata", "commands.json")
	cfg, err := Read(cfgPath, boxPath, cmdPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.PoolSize != 8 {
		t.Errorf("PoolSize: expected 8, got %d", cfg.PoolSize)
	}
	if !cfg.Verbose {
		t.Error("Verbose: expected true")
	}
	if cfg.Box.Memory != 64 {
		t.Errorf("Box.Memory: expected 64, got %d", cfg.Box.Memory)
	}
	if cfg.Step.User != "sandbox" {
		t.Errorf("Step.User: expected sandbox, got %s", cfg.Step.User)
	}
	if _, ok := cfg.Boxes["python"]; !ok {
		t.Error("Boxes: missing python box")
	}
	if _, ok := cfg.Commands["python"]; !ok {
		t.Error("Commands: missing python sandbox")
	}
	if _, ok := cfg.Commands["python"]["run"]; !ok {
		t.Error("Commands[python]: missing run command")
	}
}
