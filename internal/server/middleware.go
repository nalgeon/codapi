// HTTP middlewares.
package server

import "net/http"

// enableCORS allows cross-site requests for a given handler.
func enableCORS(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("access-control-allow-origin", "*")
		w.Header().Set("access-control-allow-method", "post")
		w.Header().Set("access-control-allow-headers", "authorization, content-type")
		w.Header().Set("access-control-max-age", "3600")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}
}
