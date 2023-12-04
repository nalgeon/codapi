package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_enableCORS(t *testing.T) {
	t.Run("options", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("OPTIONS", "/v1/exec", nil)
		handler := func(w http.ResponseWriter, r *http.Request) {}
		fn := enableCORS(handler)
		fn(w, r)

		if w.Header().Get("access-control-allow-origin") != "*" {
			t.Errorf("invalid access-control-allow-origin")
		}
		if w.Code != 200 {
			t.Errorf("expected status code 200, got %d", w.Code)
		}
	})
	t.Run("post", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/v1/exec", nil)
		handler := func(w http.ResponseWriter, r *http.Request) {}
		fn := enableCORS(handler)
		fn(w, r)

		if w.Header().Get("access-control-allow-origin") != "*" {
			t.Errorf("invalid access-control-allow-origin")
		}
		if w.Header().Get("access-control-allow-method") != "post" {
			t.Errorf("invalid access-control-allow-method")
		}
		if w.Header().Get("access-control-allow-headers") != "authorization, content-type" {
			t.Errorf("invalid access-control-allow-headers")
		}
		if w.Header().Get("access-control-max-age") != "3600" {
			t.Errorf("access-control-max-age")
		}
	})

}
