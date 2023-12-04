package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/nalgeon/codapi/internal/engine"
)

func Test_readJson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/example",
			strings.NewReader(`{"sandbox": "python", "command": "run"}`))
		req.Header.Set("Content-Type", "application/json")

		got, err := readJson[engine.Request](req)
		if err != nil {
			t.Errorf("expected nil err, got %v", err)
		}

		want := engine.Request{
			Sandbox: "python", Command: "run",
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})
	t.Run("unsupported media type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/example", nil)
		req.Header.Set("Content-Type", "text/plain")

		_, err := readJson[engine.Request](req)
		if err == nil || err.Error() != "Unsupported Media Type" {
			t.Errorf("unexpected error %v", err)
		}
	})
	t.Run("error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/example", strings.NewReader("hello world"))
		req.Header.Set("Content-Type", "application/json")

		_, err := readJson[engine.Request](req)
		if err == nil {
			t.Error("expected unmarshaling error")
		}
	})
}

func Test_writeJson(t *testing.T) {
	w := httptest.NewRecorder()
	obj := engine.Request{
		ID: "42", Sandbox: "python", Command: "run",
	}

	err := writeJson(w, obj)
	if err != nil {
		t.Errorf("expected nil err, got %v", err)
	}

	body := w.Body.String()
	contentType := w.Header().Get("content-type")
	if contentType != "application/json" {
		t.Errorf("unexpected content-type header %s", contentType)
	}

	want := `{"id":"42","sandbox":"python","command":"run","files":null}`
	if body != want {
		t.Errorf("expected %s, got %s", body, want)
	}
}

func Test_writeError(t *testing.T) {
	w := httptest.NewRecorder()
	obj := time.Date(2020, 10, 15, 0, 0, 0, 0, time.UTC)
	writeError(w, http.StatusForbidden, obj)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected status code %d, got %d", http.StatusForbidden, w.Code)
	}
	if w.Body.String() != `"2020-10-15T00:00:00Z"` {
		t.Errorf("unexpected body %s", w.Body.String())
	}
}
