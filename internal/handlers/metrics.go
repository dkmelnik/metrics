package handlers

import (
	"fmt"
	"github.com/dkmelnik/metrics/internal/handlers/interfaces"
	"net/http"
	"strings"
)

type Handler struct {
	storage interfaces.Storage
}

func NewHandler(s interfaces.Storage) *Handler {
	return &Handler{s}
}

func (h *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	metricsType := parts[2]
	metricsName := parts[3]
	metricsVal := parts[4]

	fmt.Printf("Тип метрики: %s, Имя метрики: %s, Значение метрики: %s\n", metricsType, metricsName, metricsVal)

	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, http.StatusText(http.StatusOK))
}

func (h *Handler) GetAll(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
