package metrics

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Create(rw http.ResponseWriter, r *http.Request) {
	metricsType := chi.URLParam(r, "type")
	metricsName := chi.URLParam(r, "name")
	metricsVal := chi.URLParam(r, "value")

	err := h.service.SaveMetricData(metricsType, metricsName, metricsVal)
	if err != nil {
		switch {
		case errors.Is(err, ErrTypeNotCorrect):
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		case errors.Is(err, ErrParse):
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		default:
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, http.StatusText(http.StatusOK))
}

func (h *Handler) Get(rw http.ResponseWriter, r *http.Request) {
	metricsType := chi.URLParam(r, "type")
	metricsName := chi.URLParam(r, "name")
	value, err := h.service.GetMetricData(metricsType, metricsName)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "%s", value)
}

func (h *Handler) GetAll(rw http.ResponseWriter, r *http.Request) {
	metrics := h.service.GetAllInHTML()
	b := []byte(metrics)
	_, err := rw.Write(b)
	if err != nil {
		http.Error(rw, "Failed to write response", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", http.DetectContentType(b))

}
