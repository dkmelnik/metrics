package handlers

import (
	"fmt"
	"github.com/dkmelnik/metrics/internal/handlers/interfaces"
	"net/http"
	"strconv"
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
		http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	metricsType := parts[2]
	metricsName := parts[3]
	metricsVal := parts[4]

	if metricsType == "gauge" {
		_, err := strconv.ParseFloat(metricsVal, 64)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	if metricsType == "counter" {
		_, err := strconv.ParseInt(metricsVal, 10, 64)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	fmt.Printf("Тип метрики: %s, Имя метрики: %s, Значение метрики: %f\n", metricsType, metricsName, metricsVal)

	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, http.StatusText(http.StatusOK))
}

func (h *Handler) GetAll(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
