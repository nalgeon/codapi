package config

import (
	"encoding/json"
	"os"
)

// Read reads application config from JSON files.
func Read(cfgPath, boxPath, cmdPath string) (*Config, error) {
	cfg, err := ReadConfig(cfgPath)
	if err != nil {
		return nil, err
	}

	cfg, err = ReadBoxes(cfg, boxPath)
	if err != nil {
		return nil, err
	}

	cfg, err = ReadCommands(cfg, cmdPath)
	if err != nil {
		return nil, err
	}

	return cfg, err
}

// ReadConfig reads application config from a JSON file.
func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}

// ReadBoxes reads boxes config from a JSON file.
func ReadBoxes(cfg *Config, path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	boxes := make(map[string]*Box)
	err = json.Unmarshal(data, &boxes)
	if err != nil {
		return nil, err
	}

	for _, box := range boxes {
		setBoxDefaults(box, cfg.Box)
	}

	cfg.Boxes = boxes
	return cfg, err
}

// ReadCommands reads commands config from a JSON file.
func ReadCommands(cfg *Config, path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	commands := make(map[string]SandboxCommands)
	err = json.Unmarshal(data, &commands)
	if err != nil {
		return nil, err
	}

	for _, playCmds := range commands {
		for _, cmd := range playCmds {
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

	cfg.Commands = commands
	return cfg, err
}
