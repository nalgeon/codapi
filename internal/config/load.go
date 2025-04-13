package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/nalgeon/codapi/internal/fileio"
	"github.com/nalgeon/codapi/internal/logx"
)

// Currently, Codapi supports three config layouts.
// Only the first layout is preferred, the other two will be removed in the future.
//
// 1. Sandboxes dir (preferred)
//    ├── codapi.json
//    └── sandboxes
//        ├── bash
//        │   ├── Dockerfile
//        │   ├── box.json
//        │   └── commands.json
//        └── python
//            ├── Dockerfile
//            ├── box.json
//            └── commands.json
//
// 2. Images/boxes/commands dirs (deprecated)
//    ├── configs
//    │   ├── config.json
//    │   ├── boxes
//    │   │   ├── bash.json
//    │   │   └── python.json
//    │   └── commands
//    │       ├── bash.json
//    │       └── python.json
//    └── images
//        ├── bash
//        │   └── Dockerfile
//        └── python
//            └── Dockerfile
//
// 3. Images/commands dirs + boxes.json (deprecated)
//    ├── configs
//    │   ├── config.json
//    │   ├── boxes.json
//    │   └── commands
//    │       ├── bash.json
//    │       └── python.json
//    └── images
//        ├── bash
//        │   └── Dockerfile
//        └── bash
//            └── Dockerfile

const (
	boxesDirname    = "boxes"
	codapiFilename  = "codapi.json"
	configFilename  = "config.json"
	configDirname   = "configs"
	commandsDirname = "commands"
	sandDirname     = "sandboxes"
)

var (
	ErrMissingBox  = errors.New("missing 'box' section in codapi.json")
	ErrMissingStep = errors.New("missing 'step' section in codapi.json")
)

// Read reads application config from JSON files.
func Read(path string) (*Config, error) {
	cfg, err := ReadConfig(path)
	if err != nil {
		return nil, err
	}

	cfg, err = ReadBoxes(cfg, path)
	if err != nil {
		return nil, err
	}

	cfg, err = ReadCommands(cfg, path)
	if err != nil {
		return nil, err
	}

	return cfg, err
}

// ReadConfig reads application config from a JSON file.
func ReadConfig(basePath string) (*Config, error) {
	preferredPath := filepath.Join(basePath, codapiFilename)
	fallbackPath := filepath.Join(basePath, configDirname, configFilename)

	if fileio.Exists(preferredPath) {
		return readConfig(preferredPath)
	} else {
		return readConfig(fallbackPath)
	}
}

// ReadBoxes reads boxes config from the file system.
// It prefers the sandboxes dir if it exists, otherwise fallbacks to the
// boxes dir if it exists, and finally fallbacks to the boxes.json file.
func ReadBoxes(cfg *Config, basePath string) (*Config, error) {
	var boxes map[string]*Box
	var err error

	sandDirPath := filepath.Join(basePath, sandDirname)
	boxDirPath := filepath.Join(basePath, configDirname, boxesDirname)
	boxFilePath := filepath.Join(basePath, configDirname, boxesDirname+".json")

	if fileio.Exists(sandDirPath) {
		// 1st priority is the sandboxes dir.
		boxes, err = readBoxesDir(sandDirPath, "*/box.json")
	} else if fileio.Exists(boxDirPath) {
		// 2nd priority is the configs/boxes dir.
		boxes, err = readBoxesDir(boxDirPath, "*.json")
	} else {
		// 3rd priority is configs/boxes.json.
		boxes, err = readBoxesFile(boxFilePath)
	}

	if err != nil {
		return nil, err
	}

	for _, box := range boxes {
		setBoxDefaults(box, cfg.Box)
	}

	cfg.Boxes = boxes
	return cfg, nil
}

// ReadCommands reads command configs from the file system.
// It prefers the sandboxes dir if it exists, otherwise
// fallbacks to the commands dir.
func ReadCommands(cfg *Config, basePath string) (*Config, error) {
	sandDirPath := filepath.Join(basePath, sandDirname)
	commandDirPath := filepath.Join(basePath, configDirname, commandsDirname)

	if fileio.Exists(sandDirPath) {
		// Prefer the sandboxes dir.
		return readCommands(cfg, sandDirPath, "*/commands.json")
	} else {
		// Fallback to configs/commands dir.
		return readCommands(cfg, commandDirPath, "*.json")
	}
}

// readConfig reads application config from a JSON file.
func readConfig(path string) (*Config, error) {
	logx.Debug("reading config from %s", path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Box == nil {
		return nil, ErrMissingBox
	}
	if cfg.Step == nil {
		return nil, ErrMissingStep
	}
	if cfg.HTTP == nil {
		cfg.HTTP = &HTTP{}
	}

	return cfg, err
}

// readBoxesDir reads boxes config from the boxes dir.
func readBoxesDir(path string, pattern string) (map[string]*Box, error) {
	logx.Debug("reading boxes from %s/%s", path, pattern)
	fnames, err := filepath.Glob(filepath.Join(path, pattern))
	if err != nil {
		return nil, err
	}

	boxes := make(map[string]*Box, len(fnames))
	for _, fname := range fnames {
		box, err := fileio.ReadJson[Box](fname)
		if err != nil {
			return nil, err
		}
		if box.Name == "" {
			// Determine the box name from the path.
			name := filepath.Base(fname)
			if name == "box.json" {
				// Use the parent dir name as the box name.
				name = filepath.Base(filepath.Dir(fname))
			} else {
				// Use the filename without extension as the box name.
				name = strings.TrimSuffix(name, ".json")
			}
			box.Name = name
		}
		boxes[box.Name] = &box
	}

	return boxes, err
}

// readBoxesFile reads boxes config from the boxes.json file.
func readBoxesFile(path string) (map[string]*Box, error) {
	logx.Debug("reading boxes from %s", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	boxes := make(map[string]*Box)
	err = json.Unmarshal(data, &boxes)
	if err != nil {
		return nil, err
	}

	return boxes, err
}

// readCommands reads command configs from a set of JSON files in the given path.
func readCommands(cfg *Config, path string, pattern string) (*Config, error) {
	logx.Debug("reading commands from %s/%s", path, pattern)
	fnames, err := filepath.Glob(filepath.Join(path, pattern))
	if err != nil {
		return nil, err
	}

	cfg.Commands = make(map[string]SandboxCommands, len(fnames))
	for _, fname := range fnames {
		// Determine the sandbox name from the path.
		name := filepath.Base(fname)
		if name == "commands.json" {
			// Use the parent dir name as the sandbox name.
			name = filepath.Base(filepath.Dir(fname))
		} else {
			// Use the filename without extension as the sandbox name.
			name = strings.TrimSuffix(name, ".json")
		}
		// Read the commands from the file.
		commands, err := fileio.ReadJson[SandboxCommands](fname)
		if err != nil {
			break
		}
		setCommandDefaults(commands, cfg)
		cfg.Commands[name] = commands
	}

	return cfg, err
}

// setCommandDefaults applies global defaults to sandbox commands.
func setCommandDefaults(commands SandboxCommands, cfg *Config) {
	for _, cmd := range commands {
		if cmd.Before != nil {
			setStepDefaults(cmd.Before, cfg.Step)
		}
		for _, step := range cmd.Steps {
			setStepDefaults(step, cfg.Step)
		}
		if cmd.After != nil {
			setStepDefaults(cmd.After, cfg.Step)
		}
	}
}
