package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nalgeon/be"
	"github.com/nalgeon/codapi/internal/engine"
)

func Test_readJson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/example",
			strings.NewReader(`{"sandbox": "python", "command": "run"}`))
		req.Header.Set("Content-Type", "application/json")

		got, err := readJson[engine.Request](req)
		be.Err(t, err, nil)

		want := engine.Request{
			Sandbox: "python", Command: "run",
		}
		be.Equal(t, got, want)
	})
	t.Run("unsupported media type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/example", nil)
		req.Header.Set("Content-Type", "text/plain")

		_, err := readJson[engine.Request](req)
		be.Err(t, err, "Unsupported Media Type")
	})
	t.Run("error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/example", strings.NewReader("hello world"))
		req.Header.Set("Content-Type", "application/json")

		_, err := readJson[engine.Request](req)
		be.Err(t, err)
	})
}

func Test_writeJson(t *testing.T) {
	w := httptest.NewRecorder()
	obj := engine.Request{
		ID: "42", Sandbox: "python", Command: "run",
	}

	err := writeJson(w, obj)
	be.Err(t, err, nil)

	body := w.Body.String()
	contentType := w.Header().Get("content-type")
	be.Equal(t, contentType, "application/json")

	want := `{"id":"42","sandbox":"python","command":"run","files":null}`
	be.Equal(t, body, want)
}

func Test_writeError(t *testing.T) {
	w := httptest.NewRecorder()
	obj := time.Date(2020, 10, 15, 0, 0, 0, 0, time.UTC)
	writeError(w, http.StatusForbidden, obj)
	be.Equal(t, w.Code, http.StatusForbidden)
	be.Equal(t, w.Body.String(), `"2020-10-15T00:00:00Z"`)
}
