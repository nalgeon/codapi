package config

import (
	"testing"

	"github.com/nalgeon/be"
)

func TestRead(t *testing.T) {
	cfg, err := Read("testdata")
	be.Err(t, err, nil)
	be.Equal(t, cfg.PoolSize, 8)
	be.Equal(t, cfg.Verbose, true)
	be.Equal(t, cfg.Box.Memory, 64)
	be.Equal(t, cfg.Step.User, "sandbox")

	// alpine box
	be.True(t, cfg.Boxes["custom-alpine"] != nil)
	be.Equal(t, cfg.Boxes["custom-alpine"].Image, "custom/alpine")

	// python box
	be.True(t, cfg.Boxes["python"] != nil)
	be.True(t, cfg.Commands["python"] != nil)
	be.True(t, cfg.Commands["python"]["run"] != nil)
}
