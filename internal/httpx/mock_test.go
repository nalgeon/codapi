package httpx

import (
	"io"
	"net/http"
	"testing"

	"github.com/nalgeon/be"
)

func TestMockClient(t *testing.T) {
	Mock()

	const url = "https://codapi.org/example.txt"
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := Do(req)
	be.Err(t, err, nil)
	defer func() { _ = resp.Body.Close() }()

	be.Equal(t, resp.StatusCode, http.StatusOK)

	body, err := io.ReadAll(resp.Body)
	be.Err(t, err, nil)

	want := "hello"
	be.Equal(t, string(body), want)
}
