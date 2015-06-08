//
// HTTP middleware
//

package gtcha

import (
	"fmt"
	"net/http"
)

func handle(pattern string, handler http.HandlerFunc, methods ...string) {
	http.Handle(pattern, corsMiddleware(methodMiddleware(handler, methods...)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set(
			"Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"),
		)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		next.ServeHTTP(w, r)
	})
}

func methodMiddleware(next http.Handler, methods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if len(methods) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		for _, method := range methods {
			if r.Method == method {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(
			w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed,
		)
	})
}
