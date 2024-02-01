package middlewares

import (
	"bytes"
	"encoding/hex"
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

		hash, err := hex.DecodeString(r.Header.Get("HashSHA256"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !m.sign.Equal(hash, body) {
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
