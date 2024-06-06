package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/dkmelnik/metrics/internal/apperrors"
	"github.com/dkmelnik/metrics/internal/metrics"
	"github.com/dkmelnik/metrics/internal/metrics/dto"
	"github.com/dkmelnik/metrics/internal/models"
)

// Handler is an HTTP handler for handling requests related to metrics.
type Handler struct {
	pgDB    *sqlx.DB
	service *metrics.Service
}

// NewHandler creates a new instance of Handler.
func NewHandler(pgDB *sqlx.DB, s *metrics.Service) *Handler {
	return &Handler{pgDB, s}
}

// CreateOrUpdateByParams is an HTTP handler method to create or update metrics based on URL parameters.
func (h *Handler) CreateOrUpdateByParams(rw http.ResponseWriter, r *http.Request) {
	metricsType := chi.URLParam(r, "type")
	metricsName := chi.URLParam(r, "name")
	metricsVal := chi.URLParam(r, "value")

	err := h.service.CreateOrUpdateByParams(metricsType, metricsName, metricsVal)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrTypeNotCorrect):
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		case errors.Is(err, apperrors.ErrParse):
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

// CreateOrUpdateByJSON is an HTTP handler method to create or update metrics based on JSON request body.
func (h *Handler) CreateOrUpdateByJSON(rw http.ResponseWriter, r *http.Request) {
	var body dto.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if body.Delta == nil && body.Value == nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	model, err := models.NewMetric(body.ID, body.MType)

	if err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if body.Delta != nil {
		model.SetDelta(*body.Delta)
	}
	if body.Value != nil {
		model.SetValue(*body.Value)
	}

	if err := h.service.CreateOrUpdate(model); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrTypeNotCorrect):
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		case errors.Is(err, apperrors.ErrParse):
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		default:
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var out dto.Details
	out.FillFromModel(model)

	marshal, err := json.Marshal(out)
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

// CreateOrUpdateMany is an HTTP handler method to create or update multiple metrics based on JSON request body.
func (h *Handler) CreateOrUpdateMany(rw http.ResponseWriter, r *http.Request) {
	var body []dto.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var mds = make([]models.Metric, 0, len(body))
	for _, v := range body {
		if v.Delta == nil && v.Value == nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		model, err := models.NewMetric(v.ID, v.MType)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if v.Delta != nil {
			model.SetDelta(*v.Delta)
		}
		if v.Value != nil {
			model.SetValue(*v.Value)
		}

		mds = append(mds, model)
	}

	if err := h.service.CreateOrUpdateMany(mds); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrTypeNotCorrect):
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		case errors.Is(err, apperrors.ErrParse):
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		default:
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var out = make([]dto.Details, len(mds))
	for idx, v := range mds {
		out[idx].FillFromModel(v)
	}

	marshal, err := json.Marshal(out)
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

// GetMetric is an HTTP handler method to retrieve a metric based on JSON request body.
func (h *Handler) GetMetric(rw http.ResponseWriter, r *http.Request) {
	var body dto.GetRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	value, err := h.service.GetMetric(body.MType, body.ID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
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

// GetMetricValue is an HTTP handler method to retrieve the value of a metric based on URL parameters.
func (h *Handler) GetMetricValue(rw http.ResponseWriter, r *http.Request) {
	metricsType := chi.URLParam(r, "type")
	metricsName := chi.URLParam(r, "name")
	value, err := h.service.GetMetricValue(metricsType, metricsName)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	if _, err = fmt.Fprintf(rw, "%v", value); err != nil {
		http.Error(rw, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// GetAllMetrics is an HTTP handler method to retrieve all metrics in HTML format.
func (h *Handler) GetAllMetrics(rw http.ResponseWriter, r *http.Request) {
	metrics, err := h.service.GetAllInHTML()
	if err != nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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
