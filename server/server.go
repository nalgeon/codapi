// Package server provides an HTTP API for running code in a sandbox.
package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/nalgeon/codapi/logx"
)

// The maximum duration of the server graceful shutdown.
const ShutdownTimeout = 3 * time.Second

// A Server is an HTTP sandbox server.
type Server struct {
	srv *http.Server
	wg  *sync.WaitGroup
}

// NewServer creates a new Server.
func NewServer(port int, handler http.Handler) *Server {
	addr := fmt.Sprintf(":%d", port)
	return &Server{
		srv: &http.Server{Addr: addr, Handler: handler},
		wg:  &sync.WaitGroup{},
	}
}

// Start starts the server.
func (s *Server) Start() {
	// run the server inside a goroutine so that
	// it does not block the main goroutine, and allow it
	// to start other processes and listen for signals
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := s.srv.ListenAndServe()
		if err != http.ErrServerClosed {
			logx.Log(err.Error())
		}
	}()
}

// Stop stops the server.
func (s *Server) Stop() error {
	// perform a graceful shutdown, but not longer
	// than the duration of ShutdownTimeout
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return err
	}
	s.wg.Wait()
	return nil
}
