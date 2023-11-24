package httpx

import (
	"io"
	"net/http"
	"testing"
)

func TestMockClient(t *testing.T) {
	Mock()

	const url = "https://codapi.org/example.txt"
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := Do(req)
	if err != nil {
		t.Errorf("Do: unexpected error %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Do: expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("io.ReadAll: unexpected error %v", err)
	}

	want := "hello"
	if string(body) != want {
		t.Errorf("Do: expected %v, got %v", want, string(body))
	}
}
