package server

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := NewServer(8585, handler)
	if srv.srv.Addr != ":8585" {
		t.Fatalf("NewServer: expected port :8585 got %s", srv.srv.Addr)
	}

	srv.Start()
	resp, err := http.Get("http://localhost:8585/get")
	if err != nil {
		t.Fatalf("GET: expected nil err, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET: expected status code 200, got %d", resp.StatusCode)
	}

	err = srv.Stop()
	if err != nil {
		t.Fatalf("Stop: expected nil err, got %v", err)
	}
}
