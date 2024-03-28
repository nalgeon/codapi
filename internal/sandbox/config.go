// Creates sandboxes according to the configuration.
package sandbox

import (
	"fmt"

	"github.com/nalgeon/codapi/internal/config"
	"github.com/nalgeon/codapi/internal/engine"
)

// A semaphore represents available concurrent workers
// that are responsible for executing code in sandboxes.
// The workers themselves are external to this package
// (the calling goroutines are workers).
var semaphore *Semaphore

var engineConstr = map[string]func(*config.Config, string, string) engine.Engine{
	"docker": engine.NewDocker,
	"http":   engine.NewHTTP,
}

// engines is the registry of command executors.
// Each engine executes a specific command in a specific sandbox.
// sandbox : command : engine
// TODO: Maybe it's better to create a single instance of each engine
// and pass the sandbox and command as arguments to the Exec.
var engines = map[string]map[string]engine.Engine{}

// ApplyConfig fills engine registry according to the configuration.
func ApplyConfig(cfg *config.Config) error {
	semaphore = NewSemaphore(cfg.PoolSize)
	for sandName, sandCmds := range cfg.Commands {
		engines[sandName] = make(map[string]engine.Engine)
		for cmdName, cmd := range sandCmds {
			constructor, ok := engineConstr[cmd.Engine]
			if !ok {
				return fmt.Errorf("unknown engine: %s", cmd.Engine)
			}
			engines[sandName][cmdName] = constructor(cfg, sandName, cmdName)
		}
	}
	return nil
}
