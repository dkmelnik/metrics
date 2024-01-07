package middlewares

import "net/http"

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func (m *Manager) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
