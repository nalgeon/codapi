// Reading requests and writing responses.
package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// readJson decodes the request body from JSON.
func readJson[T any](r *http.Request) (T, error) {
	var obj T
	if r.Header.Get("content-type") != "application/json" {
		return obj, errors.New(http.StatusText(http.StatusUnsupportedMediaType))
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return obj, err
	}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return obj, err
	}
	return obj, err
}

// writeJson encodes an object into JSON and writes it to the response.
func writeJson(w http.ResponseWriter, obj any) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	w.Header().Set("content-type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// writeError encodes an error object into JSON and writes it to the response.
func writeError(w http.ResponseWriter, code int, obj any) {
	data, _ := json.Marshal(obj)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data) //nolint:errcheck
}
