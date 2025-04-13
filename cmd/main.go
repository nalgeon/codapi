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

// set by the build process
var (
	version = "main"
	commit  = "none"
	date    = "unknown"
)

// startServer starts the HTTP API sandbox server.
func startServer(port int) *server.Server {
	const host = "" // listen on all interfaces
	logx.Log("codapi %s, commit %s, built at %s", version, commit, date)
	logx.Log("listening on 0.0.0.0:%d...", port)
	router := server.NewRouter()
	srv := server.NewServer(host, port, router)
	srv.Start()
	return srv
}

// startDebug servers the debug handlers.
func startDebug(port int) *server.Server {
	const host = "localhost"
	logx.Log("debugging on localhost:%d...", port)
	router := server.NewDebug()
	srv := server.NewServer(host, port, router)
	srv.Start()
	return srv
}

// listenSignals listens for termination signals
// and performs graceful shutdown.
func listenSignals(servers []*server.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	logx.Log("stopping...")
	for _, srv := range servers {
		err := srv.Stop()
		if err != nil {
			logx.Log("failed to stop: %v", err)
		}
	}
}

func main() {
	port := flag.Int("port", 1313, "server port")
	flag.Parse()

	cfg, err := config.Read(".")
	if err != nil {
		logx.Log("read config: %v", err)
		os.Exit(1)
	}

	err = sandbox.ApplyConfig(cfg)
	if err != nil {
		logx.Log("apply config: %v", err)
		os.Exit(1)
	}

	srv := startServer(*port)
	logx.Verbose = cfg.Verbose
	logx.Log("workers: %d", cfg.PoolSize)
	logx.Log("boxes: %v", cfg.BoxNames())
	logx.Log("commands: %v", cfg.CommandNames())

	debug := startDebug(6060)
	listenSignals([]*server.Server{srv, debug})
}
