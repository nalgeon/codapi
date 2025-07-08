package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nalgeon/be"
	"github.com/nalgeon/codapi/internal/config"
	"github.com/nalgeon/codapi/internal/engine"
	"github.com/nalgeon/codapi/internal/execy"
	"github.com/nalgeon/codapi/internal/sandbox"
)

var cfg = &config.Config{
	PoolSize: 8,
	Boxes: map[string]*config.Box{
		"python": {},
	},
	Commands: map[string]config.SandboxCommands{
		"python": map[string]*config.Command{
			"run": {
				Engine: "docker",
				Entry:  "main.py",
				Steps: []*config.Step{
					{Box: "python", Action: "run", NOutput: 4096},
				},
			},
			"test": {Engine: "docker"},
		},
	},
}

type server struct {
	srv *httptest.Server
	cli *http.Client
}

func newServer() *server {
	router := NewRouter()
	srv := httptest.NewServer(router)
	return &server{srv, srv.Client()}
}

func (s *server) post(uri string, val any) (*http.Response, error) {
	body, _ := json.Marshal(val)
	req, _ := http.NewRequest("POST", s.srv.URL+uri, bytes.NewReader(body))
	req.Header.Set("content-type", "application/json")
	return s.cli.Do(req)
}

func (s *server) close() {
	s.srv.Close()
}

func Test_exec(t *testing.T) {
	_ = sandbox.ApplyConfig(cfg)
	execy.Mock(map[string]execy.CmdOut{
		"docker run": {Stdout: "hello"},
	})

	srv := newServer()
	defer srv.close()

	t.Run("success", func(t *testing.T) {
		in := engine.Request{
			Sandbox: "python",
			Command: "run",
			Files: map[string]string{
				"": "print('hello')",
			},
		}
		resp, err := srv.post("/v1/exec", in)
		be.Err(t, err, nil)
		out := decodeResp[engine.Execution](t, resp)
		be.True(t, out.OK)
		be.Equal(t, out.Stdout, "hello")
		be.Equal(t, out.Stderr, "")
		be.Equal(t, out.Err, nil)
	})
	t.Run("error not found", func(t *testing.T) {
		in := engine.Request{
			Sandbox: "rust",
			Command: "run",
			Files:   nil,
		}
		resp, err := srv.post("/v1/exec", in)
		be.Err(t, err, nil)
		be.Equal(t, resp.StatusCode, http.StatusNotFound)
		out := decodeResp[engine.Execution](t, resp)
		be.Equal(t, out.OK, false)
		be.Equal(t, out.Stdout, "")
		be.Equal(t, out.Stderr, "unknown sandbox")
		be.Equal(t, out.Err, nil)
	})
	t.Run("error bad request", func(t *testing.T) {
		in := engine.Request{
			Sandbox: "python",
			Command: "run",
			Files:   nil,
		}
		resp, err := srv.post("/v1/exec", in)
		be.Err(t, err, nil)
		be.Equal(t, resp.StatusCode, http.StatusBadRequest)
		out := decodeResp[engine.Execution](t, resp)
		be.Equal(t, out.OK, false)
		be.Equal(t, out.Stdout, "")
		be.Equal(t, out.Stderr, "empty request")
		be.Equal(t, out.Err, nil)
	})
}

func decodeResp[T any](t *testing.T, resp *http.Response) T {
	defer func() { _ = resp.Body.Close() }()
	var val T
	err := json.NewDecoder(resp.Body).Decode(&val)
	be.Err(t, err, nil)
	return val
}
