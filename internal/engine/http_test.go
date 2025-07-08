package engine

import (
	"io"
	"net/http"
	"testing"

	"github.com/nalgeon/be"
	"github.com/nalgeon/codapi/internal/config"
	"github.com/nalgeon/codapi/internal/httpx"
	"github.com/nalgeon/codapi/internal/logx"
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
		be.Equal(t, out.ID, req.ID)
		be.True(t, out.OK)
		want := `HTTP/1.1 200 OK
Content-Type: text/plain

hello`
		be.Equal(t, out.Stdout, want)
		be.Equal(t, out.Stderr, "")
		be.Equal(t, out.Err, nil)
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
		be.Equal(t, out.Err, nil)
		be.Equal(t, out.Stderr, "host not allowed: example.com")
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
		be.Err(t, err, nil)
		be.Equal(t, req.Method, http.MethodHead)
		be.Equal(t, req.URL.String(), uri)
	})
	t.Run("headers", func(t *testing.T) {
		const uri = "https://codapi.org/get"
		text := "GET " + uri + "\naccept: text/plain\nx-secret: 42"
		req, err := engine.parse(text)
		be.Err(t, err, nil)
		be.Equal(t, req.Method, http.MethodGet)
		be.Equal(t, req.URL.String(), uri)
		be.Equal(t, len(req.Header), 2)
		be.Equal(t, req.Header.Get("accept"), "text/plain")
		be.Equal(t, req.Header.Get("x-secret"), "42")
	})
	t.Run("body", func(t *testing.T) {
		const uri = "https://codapi.org/post"
		const body = "{\"name\":\"alice\"}"
		text := "POST " + uri + "\ncontent-type: application/json\n\n" + body
		req, err := engine.parse(text)
		be.Err(t, err, nil)
		be.Equal(t, req.Method, http.MethodPost)
		be.Equal(t, req.URL.String(), uri)
		be.Equal(t, req.Header.Get("content-type"), "application/json")
		b, _ := io.ReadAll(req.Body)
		got := string(b)
		be.Equal(t, got, body)
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := engine.parse("on,e two three")
		be.Err(t, err)
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
		be.True(t, ok)
		be.Equal(t, req.URL.Hostname(), "localhost")
	})
	t.Run("unknown url", func(t *testing.T) {
		const uri = "http://example.com/get"
		req, _ := http.NewRequest(http.MethodGet, uri, nil)
		ok := engine.translateHost(req)
		be.Equal(t, ok, false)
		be.Equal(t, req.URL.Hostname(), "example.com")
	})
}
