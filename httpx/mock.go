package httpx

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
)

var contentTypes = map[string]string{
	".json": "application/json",
	".txt":  "text/plain",
}

// MockClient serves responses from the file system instead of remote calls.
// Should be used for testing purposes only.
type MockClient struct {
	dir string
}

// Mock creates a new MockClient and installs it instead of the default one.
func Mock(path ...string) *MockClient {
	dir := filepath.Join("testdata", filepath.Join(path...))
	c := &MockClient{dir: dir}
	client = c
	return c
}

// Do serves the file according to the request URL.
func (c *MockClient) Do(req *http.Request) (*http.Response, error) {
	filename := filepath.Join(c.dir, path.Base(req.URL.Path))

	data, err := os.ReadFile(filename)
	if err != nil {
		resp := http.Response{
			Status:     http.StatusText(http.StatusNotFound),
			StatusCode: http.StatusNotFound,
		}
		return &resp, nil
	}

	cType, ok := contentTypes[path.Ext(filename)]
	if !ok {
		cType = "application/octet-stream"
	}
	rdr := respond(cType, data)
	resp, err := http.ReadResponse(bufio.NewReader(rdr), req)
	if err != nil {
		panic(err)
	}
	return resp, nil
}

func respond(cType string, data []byte) io.Reader {
	buf := bytes.Buffer{}
	buf.WriteString("HTTP/1.1 200 OK\n")
	buf.WriteString(fmt.Sprintf("Content-Type: %s\n\n", cType))
	_, err := buf.Write(data)
	if err != nil {
		panic(err)
	}
	return &buf
}

// MockServer creates a mock HTTP server and installs its client
// instead of the default one. Serves responses from the file system
// instead of remote calls. Should be used for testing purposes only.
func MockServer() *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Join("testdata", path.Base(r.URL.Path))

		data, err := os.ReadFile(filename)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		cType, ok := contentTypes[path.Ext(filename)]
		if !ok {
			cType = "application/octet-stream"
		}

		w.Header().Set("content-type", cType)
		_, err = w.Write(data)
		if err != nil {
			panic(err)
		}
	}))
	client = srv.Client()
	return srv
}
