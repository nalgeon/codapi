package config

import (
	"path/filepath"
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

	// docker
	be.Equal(t, cfg.Docker.Bin, "docker")
	be.Equal(t, cfg.Docker.Tmp, "")

	// alpine box
	be.True(t, cfg.Boxes["custom-alpine"] != nil)
	be.Equal(t, cfg.Boxes["custom-alpine"].Image, "custom/alpine")

	// python box
	be.True(t, cfg.Boxes["python"] != nil)
	be.True(t, cfg.Commands["python"] != nil)
	be.True(t, cfg.Commands["python"]["run"] != nil)
}

func TestReadDocker(t *testing.T) {
	temp := filepath.Join(t.TempDir(), "TEMP")
	tempDir = func() string { return temp }

	cfg, err := readConfig(filepath.Join("testdata", "docker.json"))

	be.Err(t, err, nil)
	be.Equal(t, cfg.Docker.Bin, "docker")
	be.Equal(t, cfg.Docker.Tmp, temp)
}

func TestExpandTilde(t *testing.T) {
	home := filepath.Join(t.TempDir(), "HOME")
	userHomeDir = func() (string, error) { return home, nil }

	tmp, err := expandTilde("~/tmp")
	be.Err(t, err, nil)

	be.Equal(t, tmp, home+"/tmp")
}
