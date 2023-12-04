// Codapi safely executes code snippets using sandboxes.
package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/nalgeon/codapi/internal/config"
	"github.com/nalgeon/codapi/internal/logx"
	"github.com/nalgeon/codapi/internal/sandbox"
	"github.com/nalgeon/codapi/internal/server"
)

var Version string = "main"

// startServer starts the HTTP API sandbox server.
func startServer(port int) *server.Server {
	logx.Log("codapi %s", Version)
	logx.Log("listening on port %d...", port)
	router := server.NewRouter()
	srv := server.NewServer(port, router)
	srv.Start()
	return srv
}

// listenSignals listens for termination signals
// and performs graceful shutdown.
func listenSignals(srv *server.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	logx.Log("stopping...")
	err := srv.Stop()
	if err != nil {
		logx.Log("failed to stop: %v", err)
	}
}

func main() {
	port := flag.Int("port", 1313, "server port")
	flag.Parse()

	cfg, err := config.Read("config.json", "boxes.json", "commands.json")
	if err != nil {
		logx.Log("missing config file")
		os.Exit(1)
	}

	err = sandbox.ApplyConfig(cfg)
	if err != nil {
		logx.Log("invalid config: %v", err)
		os.Exit(1)
	}

	srv := startServer(*port)
	logx.Verbose = cfg.Verbose
	logx.Log("workers: %d", cfg.PoolSize)
	logx.Log("boxes: %v", cfg.BoxNames())
	logx.Log("commands: %v", cfg.CommandNames())

	listenSignals(srv)
}
