// Package config reads application config.
package config

import (
	"encoding/json"
	"sort"
)

// A Config describes application config.
type Config struct {
	PoolSize int   `json:"pool_size"`
	Verbose  bool  `json:"verbose"`
	Box      *Box  `json:"box"`
	Step     *Step `json:"step"`
	HTTP     *HTTP `json:"http"`

	// These are the available containers ("boxes").
	Boxes map[string]*Box `json:"boxes"`

	// These are the "sandboxes". Each sandbox can contain
	// multiple commands, and each command can contain
	// multiple steps. Each step is executed in a specific box.
	Commands map[string]SandboxCommands `json:"commands"`
}

// BoxNames returns configured box names.
func (cfg *Config) BoxNames() []string {
	names := make([]string, 0, len(cfg.Boxes))
	for name := range cfg.Boxes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// CommandNames returns configured command names.
func (cfg *Config) CommandNames() []string {
	names := make([]string, 0, len(cfg.Commands))
	for name := range cfg.Commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ToJSON returns JSON-encoded config with indentation.
func (cfg *Config) ToJSON() string {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return string(data)
}

// A Box describes a specific container.
// There is an important difference between a "sandbox" and a "box".
// A box is a single container. A sandbox is an environment in which we run commands.
// A sandbox command can contain multiple steps, each of which runs in a separate box.
// So the relation sandbox -> box is 1 -> 1+.
type Box struct {
	Image   string `json:"image"`
	Runtime string `json:"runtime"`
	Host

	Files []string `json:"files"`
}

// A Host describes container Host attributes.
type Host struct {
	CPU      int      `json:"cpu"`
	Memory   int      `json:"memory"`
	Storage  string   `json:"storage"`
	Network  string   `json:"network"`
	Writable bool     `json:"writable"`
	Volume   string   `json:"volume"`
	Tmpfs    []string `json:"tmpfs"`
	CapAdd   []string `json:"cap_add"`
	CapDrop  []string `json:"cap_drop"`
	Ulimit   []string `json:"ulimit"`
	// do not use the ulimit nproc because it is
	// a per-user setting, not a per-container setting
	NProc int `json:"nproc"`
}

// SandboxCommands describes all commands available for a sandbox.
// command name : command
type SandboxCommands map[string]*Command

// A Command describes a specific set of actions to take
// when executing a command in a sandbox.
type Command struct {
	Engine string  `json:"engine"`
	Entry  string  `json:"entry"`
	Before *Step   `json:"before"`
	Steps  []*Step `json:"steps"`
	After  *Step   `json:"after"`
}

// A Step describes a single step of a command.
type Step struct {
	Box     string   `json:"box"`
	Version string   `json:"version"`
	User    string   `json:"user"`
	Action  string   `json:"action"`
	Detach  bool     `json:"detach"`
	Stdin   bool     `json:"stdin"`
	Command []string `json:"command"`
	Timeout int      `json:"timeout"`
	NOutput int      `json:"noutput"`
}

// An HTTP describes HTTP engine settings.
type HTTP struct {
	Hosts map[string]string `json:"hosts"`
}

// setBoxDefaults sets default box properties
// instead of zero values.
func setBoxDefaults(box, defs *Box) {
	if box.Runtime == "" {
		box.Runtime = defs.Runtime
	}
	if box.CPU == 0 {
		box.CPU = defs.CPU
	}
	if box.Memory == 0 {
		box.Memory = defs.Memory
	}
	if box.Storage == "" {
		box.Storage = defs.Storage
	}
	if box.Network == "" {
		box.Network = defs.Network
	}
	if box.Volume == "" {
		box.Volume = defs.Volume
	}
	if box.Tmpfs == nil {
		box.Tmpfs = defs.Tmpfs
	}
	if box.CapAdd == nil {
		box.CapAdd = defs.CapAdd
	}
	if box.CapDrop == nil {
		box.CapDrop = defs.CapDrop
	}
	if box.Ulimit == nil {
		box.Ulimit = defs.Ulimit
	}
	if box.NProc == 0 {
		box.NProc = defs.NProc
	}
}

// setStepDefaults sets default command step
// properties instead of zero values.
func setStepDefaults(step, defs *Step) {
	if step.User == "" {
		step.User = defs.User
	}
	if step.Action == "" {
		step.Action = defs.Action
	}
	if step.Timeout == 0 {
		step.Timeout = defs.Timeout
	}
	if step.NOutput == 0 {
		step.NOutput = defs.NOutput
	}
}
