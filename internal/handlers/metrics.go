package handlers

import (
	"fmt"
	"github.com/dkmelnik/metrics/internal/handlers/interfaces"
	"github.com/dkmelnik/metrics/internal/models"
	"github.com/go-chi/chi/v5"
	"log"
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
	metricsType := chi.URLParam(r, "type")
	metricsName := chi.URLParam(r, "name")
	metricsVal := chi.URLParam(r, "value")

	m := models.Metrics{}

	if !m.HasType(metricsType) {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err := strconv.ParseFloat(metricsVal, 64)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = strconv.ParseInt(metricsVal, 10, 64)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.storage.Save(metricsType, metricsName, metricsVal)

	log.Printf("type:%s, name:%s, val:%s", metricsType, metricsName, metricsVal)

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, http.StatusText(http.StatusOK))
}

func (h *Handler) Get(rw http.ResponseWriter, r *http.Request) {
	metricsType := chi.URLParam(r, "type")
	metricsName := chi.URLParam(r, "name")
	metric, err := h.storage.FindOneByTypeName(metricsType, metricsName)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "%s", metric)
}

func (h *Handler) GetAll(rw http.ResponseWriter, r *http.Request) {
	metrics := h.storage.GetAllMetrics()

	html := "<html><head><title>Metric Values</title></head><body><h1>Metric Values:</h1><ul>"

	for metricName, values := range metrics {
		html += fmt.Sprintf("<li><strong>%s</strong>: <ul>", metricName)
		for key, value := range values {
			html += fmt.Sprintf("<li>%s: %v</li>", key, value)
		}
		html += "</ul></li>"
	}

	html += "</ul></body></html>"

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err := rw.Write([]byte(html))
	if err != nil {
		http.Error(rw, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
