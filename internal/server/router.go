// HTTP routes and handlers.
package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/nalgeon/codapi/internal/engine"
	"github.com/nalgeon/codapi/internal/logx"
	"github.com/nalgeon/codapi/internal/sandbox"
	"github.com/nalgeon/codapi/internal/stringx"
)

// NewRouter creates HTTP routes and handlers for them.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/exec", enableCORS(exec))
	return mux
}

// NewDebug creates HTTP routes for debugging.
func NewDebug() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return mux
}

// exec runs a sandbox command on the supplied code.
func exec(w http.ResponseWriter, r *http.Request) {
	// only POST is allowed
	if r.Method != http.MethodPost {
		err := fmt.Errorf("unsupported method: %s", r.Method)
		writeError(w, http.StatusMethodNotAllowed, engine.Fail("-", err))
		return
	}

	// read the input data - language, command, code
	in, err := readJson[engine.Request](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, engine.Fail("-", err))
		return
	}
	in.GenerateID()

	// validate the input data
	err = sandbox.Validate(in)
	if errors.Is(err, sandbox.ErrUnknownSandbox) || errors.Is(err, sandbox.ErrUnknownCommand) {
		writeError(w, http.StatusNotFound, engine.Fail(in.ID, err))
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, engine.Fail(in.ID, err))
		return
	}

	// execute the code using the sandbox
	out := sandbox.Exec(in)

	// fail on application error
	if out.Err != nil {
		logx.Log("✗ %s: %s", out.ID, out.Err)
		if errors.Is(out.Err, engine.ErrBusy) {
			writeError(w, http.StatusTooManyRequests, out)
		} else {
			writeError(w, http.StatusInternalServerError, out)
		}
		return
	}

	// log results
	if out.OK {
		logx.Log("✓ %s: took %d ms", out.ID, out.Duration)
	} else {
		msg := stringx.Compact(stringx.Shorten(out.Stderr, 80))
		logx.Log("✗ %s: %s", out.ID, msg)
	}

	// write the response
	err = writeJson(w, out)
	if err != nil {
		err = engine.NewExecutionError("write response", err)
		writeError(w, http.StatusInternalServerError, engine.Fail(in.ID, err))
		return
	}
}
