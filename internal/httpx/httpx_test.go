package httpx

import (
	"net/http"
	"testing"

	"github.com/nalgeon/be"
)

func TestDo(t *testing.T) {
	srv := MockServer()
	defer srv.Close()

	t.Run("ok", func(t *testing.T) {
		uri := srv.URL + "/example.json"
		req, _ := http.NewRequest(http.MethodGet, uri, nil)

		resp, err := Do(req)
		be.Err(t, err, nil)
		defer func() { _ = resp.Body.Close() }()

		be.Equal(t, resp.StatusCode, http.StatusOK)
	})
	t.Run("not found", func(t *testing.T) {
		uri := srv.URL + "/not-found.json"
		req, _ := http.NewRequest(http.MethodGet, uri, nil)

		resp, err := Do(req)
		be.Err(t, err, nil)
		defer func() { _ = resp.Body.Close() }()

		be.Equal(t, resp.StatusCode, http.StatusNotFound)
	})
}
