package middlewares

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"
)

func (m *MiddlewareManager) Decryption(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.privateKey == nil {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		decryptedData, err := m.decryptData(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(decryptedData))

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareManager) decryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(nil, m.privateKey, encryptedData)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}
