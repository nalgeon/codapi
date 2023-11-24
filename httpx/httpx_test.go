package httpx

import (
	"net/http"
	"testing"
)

func TestDo(t *testing.T) {
	srv := MockServer()
	defer srv.Close()

	t.Run("ok", func(t *testing.T) {
		uri := srv.URL + "/example.json"
		req, _ := http.NewRequest(http.MethodGet, uri, nil)

		resp, err := Do(req)
		if err != nil {
			t.Errorf("Do: unexpected error %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Do: expected status=%d, got %v", http.StatusOK, resp.StatusCode)
		}
	})
	t.Run("not found", func(t *testing.T) {
		uri := srv.URL + "/not-found.json"
		req, _ := http.NewRequest(http.MethodGet, uri, nil)

		resp, err := Do(req)
		if err != nil {
			t.Errorf("Do: unexpected error %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Do: expected status=%d, got %v", http.StatusNotFound, resp.StatusCode)
		}
	})
}
