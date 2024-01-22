package metrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/dkmelnik/metrics/internal/apperrors"
	"github.com/dkmelnik/metrics/internal/metrics/dto"
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

	model := models.Metric{
		Name:  body.ID,
		MType: body.MType,
	}

	if err := model.CheckType(); err != nil {
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

	var out dto.Response
	out.AdaptModel(model)

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
