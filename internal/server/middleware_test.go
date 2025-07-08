package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nalgeon/be"
)

func Test_enableCORS(t *testing.T) {
	t.Run("options", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("OPTIONS", "/v1/exec", nil)
		handler := func(w http.ResponseWriter, r *http.Request) {}
		fn := enableCORS(handler)
		fn(w, r)

		be.Equal(t, w.Header().Get("access-control-allow-origin"), "*")
		be.Equal(t, w.Code, 200)
	})
	t.Run("post", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/v1/exec", nil)
		handler := func(w http.ResponseWriter, r *http.Request) {}
		fn := enableCORS(handler)
		fn(w, r)

		be.Equal(t, w.Header().Get("access-control-allow-origin"), "*")
		be.Equal(t, w.Header().Get("access-control-allow-methods"), "options, post")
		be.Equal(t, w.Header().Get("access-control-allow-headers"), "authorization, content-type")
		be.Equal(t, w.Header().Get("access-control-max-age"), "3600")
	})

}
