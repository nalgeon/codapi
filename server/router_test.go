package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nalgeon/codapi/config"
	"github.com/nalgeon/codapi/engine"
	"github.com/nalgeon/codapi/execy"
	"github.com/nalgeon/codapi/sandbox"
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
		if err != nil {
			t.Fatalf("POST /exec: expected nil err, got %v", err)
		}
		out := decodeResp[engine.Execution](t, resp)
		if !out.OK {
			t.Error("OK: expected true")
		}
		if out.Stdout != "hello" {
			t.Errorf("Stdout: expected hello, got %s", out.Stdout)
		}
		if out.Stderr != "" {
			t.Errorf("Stderr: expected empty string, got %s", out.Stderr)
		}
		if out.Err != nil {
			t.Errorf("Err: expected nil, got %v", out.Err)
		}
	})
	t.Run("error not found", func(t *testing.T) {
		in := engine.Request{
			Sandbox: "rust",
			Command: "run",
			Files:   nil,
		}
		resp, err := srv.post("/v1/exec", in)
		if err != nil {
			t.Fatalf("POST /exec: expected nil err, got %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("StatusCode: expected 404, got %v", resp.StatusCode)
		}
		out := decodeResp[engine.Execution](t, resp)
		if out.OK {
			t.Error("OK: expected false")
		}
		if out.Stdout != "" {
			t.Errorf("Stdout: expected empty string, got %s", out.Stdout)
		}
		if out.Stderr != "unknown sandbox" {
			t.Errorf("Stderr: expected error, got %s", out.Stderr)
		}
		if out.Err != nil {
			t.Errorf("Err: expected nil, got %v", out.Err)
		}
	})
	t.Run("error bad request", func(t *testing.T) {
		in := engine.Request{
			Sandbox: "python",
			Command: "run",
			Files:   nil,
		}
		resp, err := srv.post("/v1/exec", in)
		if err != nil {
			t.Fatalf("POST /exec: expected nil err, got %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode: expected 400, got %v", resp.StatusCode)
		}
		out := decodeResp[engine.Execution](t, resp)
		if out.OK {
			t.Error("OK: expected false")
		}
		if out.Stdout != "" {
			t.Errorf("Stdout: expected empty string, got %s", out.Stdout)
		}
		if out.Stderr != "empty request" {
			t.Errorf("Stderr: expected error, got %s", out.Stderr)
		}
		if out.Err != nil {
			t.Errorf("Err: expected nil, got %v", out.Err)
		}
	})
}

func decodeResp[T any](t *testing.T, resp *http.Response) T {
	defer resp.Body.Close()
	var val T
	err := json.NewDecoder(resp.Body).Decode(&val)
	if err != nil {
		t.Fatal(err)
	}
	return val
}
