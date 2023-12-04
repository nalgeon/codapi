// Send HTTP request according to the specification.
package engine

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nalgeon/codapi/internal/config"
	"github.com/nalgeon/codapi/internal/httpx"
	"github.com/nalgeon/codapi/internal/logx"
)

// An HTTP engine sends HTTP requests.
type HTTP struct {
	hosts map[string]string
}

// NewHTTP creates a new HTTP engine.
func NewHTTP(cfg *config.Config, sandbox, command string) Engine {
	if len(cfg.HTTP.Hosts) == 0 {
		msg := fmt.Sprintf("%s %s: http engine requires at least one allowed URL", sandbox, command)
		panic(msg)
	}
	return &HTTP{hosts: cfg.HTTP.Hosts}
}

// Exec sends an HTTP request according to the spec
// and returns the response as text with status, headers and body.
func (e *HTTP) Exec(req Request) Execution {
	// build request from spec
	httpReq, err := e.parse(req.Files.First())
	if err != nil {
		err = fmt.Errorf("parse spec: %w", err)
		return Fail(req.ID, err)
	}

	// send request and receive response
	allowed := e.translateHost(httpReq)
	if !allowed {
		err = fmt.Errorf("host not allowed: %s", httpReq.Host)
		return Fail(req.ID, err)
	}

	logx.Log("%s: %s %s", req.ID, httpReq.Method, httpReq.URL.String())
	resp, err := httpx.Do(httpReq)
	if err != nil {
		err = fmt.Errorf("http request: %w", err)
		return Fail(req.ID, err)
	}
	defer resp.Body.Close()

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = NewExecutionError("read response", err)
		return Fail(req.ID, err)
	}

	// build text representation of request
	stdout := e.responseText(resp, body)
	return Execution{
		ID:     req.ID,
		OK:     true,
		Stdout: stdout,
	}
}

// parse parses the request specification.
func (e *HTTP) parse(text string) (*http.Request, error) {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return nil, errors.New("empty request")
	}

	lineIdx := 0

	// parse method and URL
	var method, url string
	methodURL := strings.Fields(lines[0])
	if len(methodURL) >= 2 {
		method = methodURL[0]
		url = methodURL[1]
	} else {
		method = http.MethodGet
		url = methodURL[0]
	}

	lineIdx++

	// parse URL parameters
	var urlParams strings.Builder
	for i := lineIdx; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "?") || strings.HasPrefix(line, "&") {
			urlParams.WriteString(line)
			lineIdx++
		} else {
			break
		}
	}

	// parse headers
	headers := make(http.Header)
	for i := lineIdx; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			break
		}
		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			headers.Add(strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1]))
			lineIdx++
		}
	}

	lineIdx += 1

	// parse body
	var bodyRdr io.Reader
	if lineIdx < len(lines) {
		body := strings.Join(lines[lineIdx:], "\n")
		bodyRdr = strings.NewReader(body)
	}

	// create request
	req, err := http.NewRequest(method, url+urlParams.String(), bodyRdr)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	return req, nil
}

// translateHost translates the requested host into the allowed one.
// Returns false if the requested host is not allowed.
func (e *HTTP) translateHost(req *http.Request) bool {
	host := e.hosts[req.Host]
	if host == "" {
		return false
	}
	req.URL.Host = host
	return true
}

// responseText returns the response as text with status, headers and body.
func (e *HTTP) responseText(resp *http.Response, body []byte) string {
	var b bytes.Buffer
	// status line
	b.WriteString(
		fmt.Sprintf("%s %d %s\n", resp.Proto, resp.StatusCode, http.StatusText(resp.StatusCode)),
	)
	// headers
	for name := range resp.Header {
		b.WriteString(fmt.Sprintf("%s: %s\n", name, resp.Header.Get(name)))
	}
	// body
	if len(body) > 0 {
		b.WriteByte('\n')
		b.Write(body)
	}
	return b.String()
}
