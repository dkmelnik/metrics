package metrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/dkmelnik/metrics/internal/metrics/dto/metric"
	"github.com/dkmelnik/metrics/internal/models"
)

type Handler struct {
	pgDB    *sqlx.DB
	service *Service
}

func NewHandler(pgDB *sqlx.DB, s *Service) *Handler {
	return &Handler{pgDB, s}
}

func (h *Handler) CreateOrUpdateByParams(rw http.ResponseWriter, r *http.Request) {
	metricsType := chi.URLParam(r, "type")
	metricsName := chi.URLParam(r, "name")
	metricsVal := chi.URLParam(r, "value")

	err := h.service.RecordMetricValue(metricsType, metricsName, metricsVal)
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

	if _, err = fmt.Fprint(rw, http.StatusText(http.StatusOK)); err != nil {
		http.Error(rw, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateOrUpdateByJSON(rw http.ResponseWriter, r *http.Request) {
	var body models.Metric
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if body.Delta == nil && body.Value == nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := h.service.ProcessMetricRequest(body); err != nil {
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

	marshal, err := json.Marshal(body)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err = rw.Write(marshal); err != nil {
		http.Error(rw, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetMetric(rw http.ResponseWriter, r *http.Request) {
	var body metric.GetRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	value, err := h.service.GetMetric(body.MType, body.ID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	marshal, err := json.Marshal(value)
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if _, err = rw.Write(marshal); err != nil {
		http.Error(rw, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetMetricValue(rw http.ResponseWriter, r *http.Request) {
	metricsType := chi.URLParam(r, "type")
	metricsName := chi.URLParam(r, "name")
	value, err := h.service.GetMetricValueString(metricsType, metricsName)
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
	if _, err = fmt.Fprintf(rw, "%s", value); err != nil {
		http.Error(rw, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAllMetrics(rw http.ResponseWriter, r *http.Request) {
	metrics := h.service.GetAllInHTML()
	b := []byte(metrics)
	rw.Header().Set("Content-Type", http.DetectContentType(b))
	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(b); err != nil {
		http.Error(rw, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CheckPostgresDBConnection(rw http.ResponseWriter, r *http.Request) {
	if h.pgDB == nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := h.pgDB.Ping(); err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}
