package engine

import (
	"io"
	"net/http"
	"testing"

	"github.com/nalgeon/codapi/config"
	"github.com/nalgeon/codapi/httpx"
	"github.com/nalgeon/codapi/logx"
)

var httpCfg = &config.Config{
	HTTP: &config.HTTP{
		Hosts: map[string]string{"codapi.org": "localhost"},
	},
	Commands: map[string]config.SandboxCommands{
		"http": map[string]*config.Command{
			"run": {Engine: "http"},
		},
	},
}

func TestHTTP_Exec(t *testing.T) {
	logx.Mock()
	httpx.Mock()
	engine := NewHTTP(httpCfg, "http", "run")

	t.Run("success", func(t *testing.T) {
		req := Request{
			ID:      "http_42",
			Sandbox: "http",
			Command: "run",
			Files: map[string]string{
				"": "GET https://codapi.org/example.txt",
			},
		}
		out := engine.Exec(req)
		if out.ID != req.ID {
			t.Errorf("ID: expected %s, got %s", req.ID, out.ID)
		}
		if !out.OK {
			t.Error("OK: expected true")
		}
		want := `HTTP/1.1 200 OK
Content-Type: text/plain

hello`
		if out.Stdout != want {
			t.Errorf("Stdout: expected %q, got %q", want, out.Stdout)
		}
		if out.Stderr != "" {
			t.Errorf("Stderr: expected %q, got %q", "", out.Stdout)
		}
		if out.Err != nil {
			t.Errorf("Err: expected nil, got %#v", out.Err)
		}
	})
	t.Run("hostname not allowed", func(t *testing.T) {
		req := Request{
			ID:      "http_42",
			Sandbox: "http",
			Command: "run",
			Files: map[string]string{
				"": "GET https://example.com/get",
			},
		}
		out := engine.Exec(req)
		if out.Err != nil {
			t.Errorf("Err: expected nil, got %#v", out.Err)
		}
		if out.Stderr != "host not allowed: example.com" {
			t.Errorf("Stderr: unexpected value %q", out.Stderr)
		}
	})
}

func TestHTTP_parse(t *testing.T) {
	logx.Mock()
	httpx.Mock()
	engine := NewHTTP(httpCfg, "http", "run").(*HTTP)

	t.Run("request line", func(t *testing.T) {
		const uri = "https://codapi.org/head"
		text := "HEAD " + uri
		req, err := engine.parse(text)
		if err != nil {
			t.Fatalf("unexpected error %#v", err)
		}
		if req.Method != http.MethodHead {
			t.Errorf("Method: expected %s, got %s", http.MethodHead, req.Method)
		}
		if req.URL.String() != uri {
			t.Errorf("URL: expected %q, got %q", uri, req.URL.String())
		}
	})
	t.Run("headers", func(t *testing.T) {
		const uri = "https://codapi.org/get"
		text := "GET " + uri + "\naccept: text/plain\nx-secret: 42"
		req, err := engine.parse(text)
		if err != nil {
			t.Fatalf("unexpected error %#v", err)
		}
		if req.Method != http.MethodGet {
			t.Errorf("Method: expected %s, got %s", http.MethodGet, req.Method)
		}
		if req.URL.String() != uri {
			t.Errorf("URL: expected %q, got %q", uri, req.URL.String())
		}
		if len(req.Header) != 2 {
			t.Fatalf("Header: expected 2 headers, got %d", len(req.Header))
		}
		if req.Header.Get("accept") != "text/plain" {
			t.Fatalf("Header: expected accept = %q, got %q", "text/plain", req.Header.Get("accept"))
		}
		if req.Header.Get("x-secret") != "42" {
			t.Fatalf("Header: expected x-secret = %q, got %q", "42", req.Header.Get("x-secret"))
		}
	})
	t.Run("body", func(t *testing.T) {
		const uri = "https://codapi.org/post"
		const body = "{\"name\":\"alice\"}"
		text := "POST " + uri + "\ncontent-type: application/json\n\n" + body
		req, err := engine.parse(text)
		if err != nil {
			t.Fatalf("unexpected error %#v", err)
		}
		if req.Method != http.MethodPost {
			t.Errorf("Method: expected %s, got %s", http.MethodPost, req.Method)
		}
		if req.URL.String() != uri {
			t.Errorf("URL: expected %q, got %q", uri, req.URL.String())
		}
		if req.Header.Get("content-type") != "application/json" {
			t.Errorf("Header: expected content-type = %q, got %q",
				"application/json", req.Header.Get("content-type"))
		}
		b, _ := io.ReadAll(req.Body)
		got := string(b)
		if got != body {
			t.Errorf("Body: expected %q, got %q", body, got)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := engine.parse("on,e two three")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestHTTP_translateHost(t *testing.T) {
	logx.Mock()
	httpx.Mock()
	engine := NewHTTP(httpCfg, "http", "run").(*HTTP)

	t.Run("known url", func(t *testing.T) {
		const uri = "http://codapi.org/get"
		req, _ := http.NewRequest(http.MethodGet, uri, nil)
		ok := engine.translateHost(req)
		if !ok {
			t.Errorf("%s: should be allowed", uri)
		}
		if req.URL.Hostname() != "localhost" {
			t.Errorf("%s: expected %s, got %s", uri, "localhost", req.URL.Hostname())
		}
	})
	t.Run("unknown url", func(t *testing.T) {
		const uri = "http://example.com/get"
		req, _ := http.NewRequest(http.MethodGet, uri, nil)
		ok := engine.translateHost(req)
		if ok {
			t.Errorf("%s: should not be allowed", uri)
		}
		if req.URL.Hostname() != "example.com" {
			t.Errorf("%s: expected %s, got %s", uri, "example.com", req.URL.Hostname())
		}
	})
}
