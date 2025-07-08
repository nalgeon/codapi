package server

import (
	"net/http"
	"testing"

	"github.com/nalgeon/be"
)

func TestServer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := NewServer("", 8585, handler)
	be.Equal(t, srv.srv.Addr, ":8585")

	srv.Start()
	resp, err := http.Get("http://localhost:8585/get")
	be.Err(t, err, nil)
	be.Equal(t, resp.StatusCode, http.StatusOK)

	err = srv.Stop()
	be.Err(t, err, nil)
}
