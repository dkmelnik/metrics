package middlewares

import (
	"fmt"
	"github.com/dkmelnik/metrics/internal/logger"
	"net/http"
)

func (m *MiddlewareManager) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				logger.Log.ErrorWithContext(r.Context(), fmt.Errorf("recover:%v", err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
