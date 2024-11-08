package config

import (
	"testing"
)

func TestRead(t *testing.T) {
	cfg, err := Read("testdata")
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

	// alpine box
	if _, ok := cfg.Boxes["custom-alpine"]; !ok {
		t.Error("Boxes: missing my/alpine box")
	}
	if cfg.Boxes["custom-alpine"].Image != "custom/alpine" {
		t.Errorf(
			"Boxes[custom-alpine]: expected custom/alpine image, got %s",
			cfg.Boxes["custom-alpine"].Image,
		)
	}

	// python box
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
