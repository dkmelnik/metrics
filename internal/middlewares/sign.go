package middlewares

import (
	"bytes"
	"io"
	"net/http"
)

func (m *MiddlewareManager) Sign(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.sign == nil {
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("HashSHA256") == "" {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		if !m.sign.Equal(r.Header.Get("HashSHA256"), body) {
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
